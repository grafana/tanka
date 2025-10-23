use anyhow::{anyhow, Context, Result};
use serde_json::Value;
use std::path::{Path, PathBuf};

use crate::environment::{Environment, LoadedEnvironment};
use crate::jsonnet::{
    compute_import_paths, find_environments_recursive, JsonnetEvaluator, SINGLE_ENV_EVAL_SCRIPT,
};
use std::fs;

/// LoaderOpts contains options for loading environments
#[derive(Debug, Clone, Default)]
pub struct LoaderOpts {
    pub name: Option<String>,
    pub cache_path: Option<PathBuf>,
}

/// Loader trait defines the interface for loading environments
pub trait Loader {
    /// Name of the loader
    fn name(&self) -> &str;

    /// Load a single environment at path
    fn load(&self, path: &Path, opts: &LoaderOpts) -> Result<LoadedEnvironment>;

    /// Peek only loads metadata and omits the actual resources
    fn peek(&self, path: &Path, opts: &LoaderOpts) -> Result<Environment>;

    /// List returns metadata of all possible environments at path that can be loaded
    fn list(&self, path: &Path, opts: &LoaderOpts) -> Result<Vec<Environment>>;

    /// Eval returns the raw evaluated Jsonnet
    fn eval(&self, path: &Path, opts: &LoaderOpts) -> Result<Value>;
}

/// InlineLoader loads an environment that is specified inline from within Jsonnet.
/// The Jsonnet output is expected to hold a tanka.dev/Environment type,
/// Kubernetes resources are expected at the `data` key of this very type
pub struct InlineLoader {}

impl InlineLoader {
    pub fn new() -> Self {
        InlineLoader {}
    }

    fn load_impl(
        &self,
        path: &Path,
        opts: &LoaderOpts,
        eval_script: Option<&str>,
    ) -> Result<LoadedEnvironment> {
        let import_paths = compute_import_paths(path);

        // Look for common Tanka files
        let main_file = if path.is_dir() {
            let main_jsonnet = path.join("main.jsonnet");
            if main_jsonnet.exists() {
                main_jsonnet
            } else {
                return Err(anyhow!("No main.jsonnet found in directory: {:?}", path));
            }
        } else {
            path.to_path_buf()
        };

        // Create a new evaluator each time, just like jrsonnet CLI
        let evaluator = JsonnetEvaluator::new_with_paths(import_paths)?;
        let value = if let Some(script) = eval_script {
            evaluator.eval_script(&main_file, script)?
        } else {
            evaluator.eval_file(&main_file)?
        };

        // Extract environments from the value
        let environments = extract_envs(&value)?;

        // Handle filtering by name
        let mut filtered_envs = environments.clone();
        if let Some(name) = &opts.name {
            // Filter environments that match the name
            filtered_envs = environments
                .iter()
                .filter(|e| {
                    e.metadata
                        .name
                        .as_ref()
                        .map(|n| n.contains(name))
                        .unwrap_or(false)
                })
                .cloned()
                .collect();

            // If there's a full match, use only that one
            if let Some(exact_match) = environments
                .iter()
                .find(|e| e.metadata.name.as_ref().map(|n| n == name).unwrap_or(false))
            {
                filtered_envs = vec![exact_match.clone()];
            }
        }

        if filtered_envs.len() > 1 {
            let names: Vec<String> = filtered_envs
                .iter()
                .filter_map(|e| e.metadata.name.clone())
                .collect();
            let name_filter = opts.name.as_ref().map(|s| s.as_str()).unwrap_or("");
            if name_filter.is_empty() {
                return Err(anyhow!(
                    "found multiple Environments in {:?}. Use `--name` to select a single one: \n - {}",
                    path,
                    names.join("\n - ")
                ));
            } else {
                return Err(anyhow!(
                    "found multiple Environments in {:?} matching {:?}. Provide a more specific name that matches a single one: \n - {}",
                    path,
                    name_filter,
                    names.join("\n - ")
                ));
            }
        }

        if filtered_envs.is_empty() {
            return Err(anyhow!(
                "found no matching environments; run 'tk env list {}' to view available options",
                path.display()
            ));
        }

        Ok(LoadedEnvironment {
            path: main_file,
            environment: filtered_envs.into_iter().next().unwrap(),
        })
    }
}

impl Loader for InlineLoader {
    fn name(&self) -> &str {
        "inline"
    }

    fn load(&self, path: &Path, opts: &LoaderOpts) -> Result<LoadedEnvironment> {
        let eval_script = opts
            .name
            .as_ref()
            .map(|n| SINGLE_ENV_EVAL_SCRIPT.replace("%s", n));
        self.load_impl(path, opts, eval_script.as_deref())
    }

    fn peek(&self, path: &Path, opts: &LoaderOpts) -> Result<Environment> {
        // Load the full environment and strip the data
        let mut loaded = self.load(path, opts)?;
        loaded.environment.data = None;
        Ok(loaded.environment)
    }

    fn list(&self, path: &Path, _opts: &LoaderOpts) -> Result<Vec<Environment>> {
        let import_paths = compute_import_paths(path);

        let main_file = if path.is_dir() {
            path.join("main.jsonnet")
        } else {
            path.to_path_buf()
        };

        // Create a new evaluator each time, just like jrsonnet CLI
        let evaluator = JsonnetEvaluator::new_with_paths(import_paths)?;
        let value = evaluator.eval_file(&main_file)?;

        // Extract environments and strip their data
        let mut envs = extract_envs(&value)?;
        for env in &mut envs {
            env.data = None;
        }
        Ok(envs)
    }

    fn eval(&self, path: &Path, _opts: &LoaderOpts) -> Result<Value> {
        let import_paths = compute_import_paths(path);

        let main_file = if path.is_dir() {
            path.join("main.jsonnet")
        } else {
            path.to_path_buf()
        };

        // Create a new evaluator each time, just like jrsonnet CLI
        let evaluator = JsonnetEvaluator::new_with_paths(import_paths)?;
        evaluator.eval_file(&main_file)
    }
}

/// StaticLoader loads an environment from a static file called `spec.json`.
/// Jsonnet is evaluated as normal
pub struct StaticLoader {}

impl StaticLoader {
    pub fn new() -> Self {
        StaticLoader {}
    }

    fn parse_static_spec(path: &Path) -> Result<Environment> {
        let spec_file = path.join("spec.json");
        let spec_contents = fs::read_to_string(&spec_file)
            .with_context(|| format!("Failed to read spec.json: {:?}", spec_file))?;

        let mut environment: Environment = serde_json::from_str(&spec_contents)
            .with_context(|| format!("Failed to parse spec.json: {:?}", spec_file))?;

        // Populate metadata.namespace from spec.namespace if not set
        if environment.metadata.namespace.is_none() {
            environment.metadata.namespace = Some(environment.spec.namespace.clone());
        }

        Ok(environment)
    }
}

impl Loader for StaticLoader {
    fn name(&self) -> &str {
        "static"
    }

    fn load(&self, path: &Path, opts: &LoaderOpts) -> Result<LoadedEnvironment> {
        let mut environment = Self::parse_static_spec(path)?;
        let data = self.eval(path, opts)?;
        environment.data = Some(data);

        let main_file = path.join("main.jsonnet");
        Ok(LoadedEnvironment {
            path: main_file,
            environment,
        })
    }

    fn peek(&self, path: &Path, _opts: &LoaderOpts) -> Result<Environment> {
        Self::parse_static_spec(path)
    }

    fn list(&self, path: &Path, opts: &LoaderOpts) -> Result<Vec<Environment>> {
        let env = self.peek(path, opts)?;
        Ok(vec![env])
    }

    fn eval(&self, path: &Path, _opts: &LoaderOpts) -> Result<Value> {
        let import_paths = compute_import_paths(path);
        let main_file = path.join("main.jsonnet");

        if !main_file.exists() {
            return Err(anyhow!("No main.jsonnet found in directory: {:?}", path));
        }

        // Create a new evaluator each time, just like jrsonnet CLI
        let evaluator = JsonnetEvaluator::new_with_paths(import_paths)?;
        evaluator.eval_file(&main_file)
    }
}

/// DetectLoader detects whether the environment is inline or static and picks
/// the appropriate loader
pub fn detect_loader(path: &Path) -> Result<Box<dyn Loader>> {
    let env_dir = if path.is_dir() {
        path
    } else {
        path.parent().unwrap_or(path)
    };

    let spec_file = env_dir.join("spec.json");
    if spec_file.exists() {
        Ok(Box::new(StaticLoader::new()))
    } else {
        Ok(Box::new(InlineLoader::new()))
    }
}

/// extractEnvs filters out any Environment manifests
fn extract_envs(data: &Value) -> Result<Vec<Environment>> {
    // Find all environments recursively in the structure
    let environments = find_environments_recursive(data)?;

    if environments.is_empty() {
        return Err(anyhow!("No Environment objects found"));
    }

    Ok(environments)
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_load_static() {
        let base_dir = PathBuf::from("testdata/test-export-envs/static-env");
        let loader = StaticLoader::new();

        let result = loader.load(&base_dir, &LoaderOpts::default());
        assert!(
            result.is_ok(),
            "Failed to load static environment: {:?}",
            result.err()
        );

        let loaded = result.unwrap();
        assert_eq!(loaded.environment.metadata.name, Some("static".to_string()));
        assert!(loaded.environment.data.is_some());
    }

    #[test]
    fn test_load_inline() {
        let base_dir = PathBuf::from("testdata/test-export-envs/inline-envs");
        let loader = InlineLoader::new();

        let opts = LoaderOpts {
            name: Some("inline-namespace1".to_string()),
            ..Default::default()
        };
        let result = loader.load(&base_dir, &opts);
        assert!(
            result.is_ok(),
            "Failed to load inline environment: {:?}",
            result.err()
        );

        let loaded = result.unwrap();
        assert_eq!(
            loaded.environment.metadata.name,
            Some("inline-namespace1".to_string())
        );
        assert!(loaded.environment.data.is_some());
    }

    #[test]
    fn test_detect_loader_static() {
        let base_dir = PathBuf::from("testdata/test-export-envs/static-env");
        let loader = detect_loader(&base_dir);
        assert!(loader.is_ok());
        assert_eq!(loader.unwrap().name(), "static");
    }

    #[test]
    fn test_detect_loader_inline() {
        let base_dir = PathBuf::from("testdata/test-export-envs/inline-envs");
        let loader = detect_loader(&base_dir);
        assert!(loader.is_ok());
        assert_eq!(loader.unwrap().name(), "inline");
    }
}
