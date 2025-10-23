use anyhow::{anyhow, Context, Result};
use log::{info, warn};
use rayon::prelude::*;
use regex::Regex;
use std::collections::HashMap;
use std::fs;
use std::io::ErrorKind;
use std::path::{Path, PathBuf};
use walkdir::WalkDir;

use crate::environment::LoadedEnvironment;
use crate::jsonnet::{load_environment, load_manifests, list_environments};
use crate::template::TemplateEngine;

const MANIFEST_FILE: &str = "manifest.json";

#[derive(Debug, Clone)]
pub enum ExportMergeStrategy {
    None,
    FailOnConflicts,
    ReplaceEnvs,
}

#[derive(Debug, Clone)]
pub struct ExportOptions {
    pub format: String,
    pub extension: String,
    pub parallelism: usize,
    pub cache_path: Option<PathBuf>,
    pub cache_path_regexes: Vec<Regex>,
    pub merge_strategy: ExportMergeStrategy,
    pub merge_deleted_envs: Vec<String>,
    pub skip_manifest: bool,
    pub targets: Vec<String>,
    pub name: Option<String>,
    pub selector: Option<String>,
}

pub async fn export_environments(
    output_dir: PathBuf,
    paths: Vec<PathBuf>,
    recursive: bool,
    opts: ExportOptions,
) -> Result<()> {
    // Check if directory is empty
    let empty = is_dir_empty(&output_dir)?;
    if !empty && matches!(opts.merge_strategy, ExportMergeStrategy::None) {
        return Err(anyhow!(
            "output dir '{}' not empty. Pass a different --merge-strategy to ignore this",
            output_dir.display()
        ));
    }

    // Find environments
    let environments = if recursive {
        find_environments_recursive(&paths, &opts)?
    } else {
        if paths.len() > 1 {
            return Err(anyhow!(
                "recursive flag is required when exporting multiple environments"
            ));
        }
        vec![load_environment(&paths[0])?]
    };

    info!("Found {} environment(s)", environments.len());

    // Delete previously exported manifests if using replace strategy
    if matches!(opts.merge_strategy, ExportMergeStrategy::ReplaceEnvs) {
        delete_previously_exported_manifests_from_envs(&output_dir, &environments, opts.skip_manifest)?;
    }

    // Delete manifests from deleted environments
    if !opts.merge_deleted_envs.is_empty() {
        delete_previously_exported_manifests(&output_dir, &opts.merge_deleted_envs, opts.skip_manifest)?;
    }

    // Export environments in parallel
    let file_to_env: HashMap<String, String> = environments
        .par_iter()
        .map(|env| export_single_environment(env, &output_dir, &opts))
        .collect::<Result<Vec<_>>>()?
        .into_iter()
        .flatten()
        .collect();

    // Write manifest file
    if !opts.skip_manifest {
        export_manifest_file(&output_dir, &file_to_env, &[])?;
    }

    info!("Successfully exported {} environment(s)", environments.len());
    Ok(())
}

fn find_environments_recursive(paths: &[PathBuf], opts: &ExportOptions) -> Result<Vec<LoadedEnvironment>> {
    let mut environments = Vec::new();

    for path in paths {
        for entry in WalkDir::new(path)
            .follow_links(true)
            .into_iter()
            .filter_map(|e| e.ok())
        {
            if entry.file_name() == "main.jsonnet" {
                // Use list_environments which can handle files with multiple environments
                match list_environments(entry.path()) {
                    Ok(envs) => {
                        for env in envs {
                            // Filter by name if specified
                            if let Some(ref name) = opts.name {
                                if let Some(ref env_name) = env.environment.metadata.name {
                                    if env_name != name {
                                        continue;
                                    }
                                } else {
                                    continue;
                                }
                            }

                            // TODO: Filter by selector if specified
                            environments.push(env);
                        }
                    }
                    Err(e) => {
                        warn!("Failed to load environment at {:?}: {}", entry.path(), e);
                    }
                }
            }
        }
    }

    Ok(environments)
}

fn export_single_environment(
    env: &LoadedEnvironment,
    output_dir: &Path,
    opts: &ExportOptions,
) -> Result<HashMap<String, String>> {
    info!("Exporting environment: {:?}", env.path);

    let manifests = load_manifests(env)?;
    
    if manifests.is_empty() {
        info!("No manifests found in environment: {:?}", env.path);
        return Ok(HashMap::new());
    }

    let template = TemplateEngine::new(&opts.format, &env.environment)?;
    let mut file_to_env = HashMap::new();

    for manifest in manifests {
        let name = template.apply(&manifest, &env.environment)?;
        let relpath = format!("{}.{}", name, opts.extension);
        let path = output_dir.join(&relpath);

        // Check if file exists (for fail-on-conflicts mode)
        if path.exists() && matches!(opts.merge_strategy, ExportMergeStrategy::FailOnConflicts) {
            return Err(anyhow!("file '{}' already exists. Aborting", path.display()));
        }

        // Get environment namespace
        let env_namespace = env.environment.metadata.namespace.as_deref().unwrap_or("default");
        file_to_env.insert(relpath, env_namespace.to_string());

        // Write manifest
        let yaml = manifest.to_yaml()?;
        write_export_file(&path, yaml.as_bytes())?;
    }

    Ok(file_to_env)
}

fn is_dir_empty(path: &Path) -> Result<bool> {
    match fs::read_dir(path) {
        Ok(mut entries) => Ok(entries.next().is_none()),
        Err(e) if e.kind() == ErrorKind::NotFound => {
            fs::create_dir_all(path)?;
            Ok(true)
        }
        Err(e) => Err(e.into()),
    }
}

fn write_export_file(path: &Path, data: &[u8]) -> Result<()> {
    if let Some(parent) = path.parent() {
        fs::create_dir_all(parent)
            .with_context(|| format!("Failed to create directory: {:?}", parent))?;
    }

    fs::write(path, data)
        .with_context(|| format!("Failed to write file: {:?}", path))?;

    Ok(())
}

fn delete_previously_exported_manifests_from_envs(
    path: &Path,
    envs: &[LoadedEnvironment],
    skip_manifest: bool,
) -> Result<()> {
    let env_names: Vec<String> = envs
        .iter()
        .filter_map(|e| e.environment.metadata.namespace.clone())
        .collect();
    
    delete_previously_exported_manifests(path, &env_names, skip_manifest)
}

fn delete_previously_exported_manifests(
    path: &Path,
    tanka_env_names: &[String],
    skip_manifest: bool,
) -> Result<()> {
    if tanka_env_names.is_empty() {
        return Ok(());
    }

    let env_names_set: HashMap<&str, ()> = tanka_env_names
        .iter()
        .map(|s| (s.as_str(), ()))
        .collect();

    let manifest_file_path = path.join(MANIFEST_FILE);
    let manifest_content = match fs::read_to_string(&manifest_file_path) {
        Ok(content) => content,
        Err(e) if e.kind() == ErrorKind::NotFound => {
            warn!(
                "No manifest file found at {:?}, skipping deletion of previously exported manifests",
                manifest_file_path
            );
            return Ok(());
        }
        Err(e) => return Err(e.into()),
    };

    let file_to_env_map: HashMap<String, String> = serde_json::from_str(&manifest_content)
        .context("Failed to parse manifest file")?;

    let mut deleted_manifest_keys = Vec::new();
    for (exported_manifest, manifest_env) in &file_to_env_map {
        if env_names_set.contains_key(manifest_env.as_str()) {
            deleted_manifest_keys.push(exported_manifest.clone());
            let file_path = path.join(exported_manifest);
            if let Err(e) = fs::remove_file(&file_path) {
                if e.kind() != ErrorKind::NotFound {
                    return Err(e.into());
                }
            }
        }
    }

    // Update manifest file
    if !skip_manifest {
        export_manifest_file(path, &HashMap::new(), &deleted_manifest_keys)?;
    }

    Ok(())
}

fn export_manifest_file(
    path: &Path,
    new_file_to_env_map: &HashMap<String, String>,
    deleted_keys: &[String],
) -> Result<()> {
    if new_file_to_env_map.is_empty() && deleted_keys.is_empty() {
        return Ok(());
    }

    let manifest_file_path = path.join(MANIFEST_FILE);
    let mut current_file_to_env_map: HashMap<String, String> = match fs::read_to_string(&manifest_file_path) {
        Ok(content) => serde_json::from_str(&content)?,
        Err(e) if e.kind() == ErrorKind::NotFound => HashMap::new(),
        Err(e) => return Err(e.into()),
    };

    // Merge new entries
    for (k, v) in new_file_to_env_map {
        current_file_to_env_map.insert(k.clone(), v.clone());
    }

    // Delete removed entries
    for k in deleted_keys {
        current_file_to_env_map.remove(k);
    }

    // Write manifest file
    let data = serde_json::to_string_pretty(&current_file_to_env_map)?;
    write_export_file(&manifest_file_path, data.as_bytes())?;

    Ok(())
}

