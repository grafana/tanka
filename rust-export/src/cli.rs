use anyhow::{anyhow, Result};
use clap::Parser;
use regex::Regex;
use std::path::PathBuf;

use crate::export::{ExportMergeStrategy, ExportOptions};

#[derive(Parser, Debug)]
#[command(name = "tanka-export")]
#[command(about = "Export Tanka environments to YAML files", long_about = None)]
pub struct Cli {
    /// Output directory for exported manifests
    #[arg(value_name = "OUTPUT_DIR")]
    pub output_dir: PathBuf,

    /// Path(s) to Tanka environments
    #[arg(value_name = "PATH", required = true)]
    pub paths: Vec<PathBuf>,

    /// Format string for output filenames
    /// See: https://tanka.dev/exporting#filenames
    #[arg(
        long,
        default_value = "{{.apiVersion}}.{{.kind}}-{{or .metadata.name .metadata.generateName}}"
    )]
    pub format: String,

    /// File extension for exported files
    #[arg(long, default_value = "yaml")]
    pub extension: String,

    /// Number of environments to process in parallel
    #[arg(short = 'p', long, default_value = "8")]
    pub parallel: usize,

    /// Local file path where cached evaluations should be stored
    #[arg(short = 'c', long)]
    pub cache_path: Option<PathBuf>,

    /// Regexes which define which environment should be cached
    #[arg(short = 'e', long = "cache-envs")]
    pub cache_envs: Vec<String>,

    /// Size of memory ballast to allocate (bytes)
    #[arg(long, default_value = "0")]
    pub mem_ballast_size_bytes: usize,

    /// What to do when exporting to an existing directory
    /// Values: 'fail-on-conflicts', 'replace-envs'
    #[arg(long)]
    pub merge_strategy: Option<String>,

    /// Tanka main files that have been deleted
    #[arg(long)]
    pub merge_deleted_envs: Vec<String>,

    /// Skip generating manifest.json file
    #[arg(long)]
    pub skip_manifest: bool,

    /// Look recursively for Tanka environments
    #[arg(short = 'r', long)]
    pub recursive: bool,

    /// Filter by resource kind/name (e.g., deployment/my-app)
    #[arg(long)]
    pub targets: Vec<String>,

    /// Filter environments by name
    #[arg(long)]
    pub name: Option<String>,

    /// Label selector to filter environments
    #[arg(long)]
    pub selector: Option<String>,
}

impl Cli {
    pub async fn run(self) -> Result<()> {
        // Validate arguments
        if !self.recursive && self.paths.len() > 1 {
            return Err(anyhow!(
                "recursive flag is required when exporting multiple environments"
            ));
        }

        // Compile cache path regexes
        let mut cache_path_regexes = Vec::new();
        for expr in &self.cache_envs {
            let regex = Regex::new(expr)?;
            cache_path_regexes.push(regex);
        }

        // Determine merge strategy
        let merge_strategy = match self.merge_strategy.as_deref() {
            Some("fail-on-conflicts") => ExportMergeStrategy::FailOnConflicts,
            Some("replace-envs") => ExportMergeStrategy::ReplaceEnvs,
            Some(s) => return Err(anyhow!("invalid merge strategy: {}", s)),
            None => ExportMergeStrategy::None,
        };

        let opts = ExportOptions {
            format: self.format,
            extension: self.extension,
            parallelism: self.parallel,
            cache_path: self.cache_path,
            cache_path_regexes,
            merge_strategy,
            merge_deleted_envs: self.merge_deleted_envs,
            skip_manifest: self.skip_manifest,
            targets: self.targets,
            name: self.name,
            selector: self.selector,
        };

        // Execute export
        crate::export::export_environments(self.output_dir, self.paths, self.recursive, opts).await
    }
}
