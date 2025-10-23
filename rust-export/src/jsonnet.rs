use anyhow::{anyhow, Context, Result};
use jrsonnet_evaluator::manifest::{JsonFormat, ManifestFormat};
use jrsonnet_evaluator::trace::PathResolver;
use jrsonnet_evaluator::{FileImportResolver, State, Val};
use jrsonnet_stdlib::ContextInitializer;
use serde_json::Value;
use std::cell::RefCell;
use std::collections::HashMap;
use std::fs;
use std::path::{Path, PathBuf};

use crate::environment::{Environment, LoadedEnvironment};
use crate::manifest::Manifest;

// Script to extract environment metadata without evaluating the data field
// This is much faster when we only need environment metadata (e.g., for listing)
const METADATA_EVAL_SCRIPT: &str = r#"
local noDataEnv(object) =
  std.prune(
    if std.isObject(object)
    then
      if std.objectHas(object, 'apiVersion')
         && std.objectHas(object, 'kind')
      then
        if object.kind == 'Environment'
        then object { data:: {} }
        else {}
      else
        std.mapWithKey(
          function(key, obj)
            noDataEnv(obj),
          object
        )
    else if std.isArray(object)
    then
      std.map(
        function(obj)
          noDataEnv(obj),
        object
      )
    else {}
  );

noDataEnv(main)
"#;

// Script to load a single environment by name with full data
// %s will be replaced with the environment name
const SINGLE_ENV_EVAL_SCRIPT: &str = r#"
local singleEnv(object) =
  if std.isObject(object)
  then
    if std.objectHas(object, 'apiVersion')
       && std.objectHas(object, 'kind')
    then
      if object.kind == 'Environment'
      && std.member(object.metadata.name, '%s')
      then object
      else {}
    else
      std.mapWithKey(
        function(key, obj)
          singleEnv(obj),
        object
      )
  else if std.isArray(object)
  then
    std.map(
      function(obj)
        singleEnv(obj),
      object
    )
  else {};

singleEnv(main)
"#;

pub struct JsonnetEvaluator {
    state: State,
}

impl JsonnetEvaluator {
    pub fn new() -> Result<Self> {
        let state = State::default();
        Ok(JsonnetEvaluator { state })
    }

    pub fn new_with_paths(paths: Vec<PathBuf>) -> Result<Self> {
        use log::debug;

        debug!("Creating evaluator with {} import paths", paths.len());
        for path in &paths {
            debug!("  - {:?}", path);
        }

        // Create FileImportResolver with paths (as done in jrsonnet CLI)
        // Filter out empty paths
        let library_paths: Vec<PathBuf> = paths
            .into_iter()
            .filter(|p| !p.as_os_str().is_empty())
            .collect();
        let import_resolver = FileImportResolver::new(library_paths);

        // Create context initializer with standard library
        let context_init = ContextInitializer::new(PathResolver::new_cwd_fallback());

        // Add Tanka native functions to the context
        crate::native::register_native_functions(&context_init);

        // Build state with both import resolver and stdlib
        let mut builder = State::builder();
        builder
            .import_resolver(import_resolver)
            .context_initializer(context_init);
        let state = builder.build();

        Ok(JsonnetEvaluator { state })
    }

    pub fn with_import_paths(&mut self, _paths: Vec<PathBuf>) -> &mut Self {
        // Can't modify import paths after State creation in this API version
        // This is kept for API compatibility but does nothing
        self
    }

    pub fn eval_file(&self, path: &Path) -> Result<Value> {
        let code = fs::read_to_string(path)
            .map_err(|e| anyhow!("Failed to read file {:?}: {}", path, e))?;

        let name = path.to_string_lossy().to_string();

        let val = self
            .state
            .evaluate_snippet(name, code)
            .map_err(|e| anyhow!("Failed to evaluate Jsonnet file {:?}: {}", path, e))?;

        self.val_to_json(&val)
    }

    pub fn eval_script(&self, path: &Path, script: &str) -> Result<Value> {
        // Read the main file to set up imports
        let main_path = path.to_string_lossy().to_string();

        // Build a script that imports the main file and runs the eval script
        let wrapped_script = format!(
            "local main = import '{}'; {}",
            main_path.replace('\\', "\\\\").replace('\'', "\\'"),
            script
        );

        let val = self
            .state
            .evaluate_snippet("<eval_script>".to_string(), wrapped_script)
            .map_err(|e| anyhow!("Failed to evaluate Jsonnet script for {:?}: {}", path, e))?;

        self.val_to_json(&val)
    }

    fn val_to_json(&self, val: &Val) -> Result<Value> {
        let json_format = JsonFormat::default();
        let mut json_str = String::new();
        json_format
            .manifest_buf(val.clone(), &mut json_str)
            .map_err(|e| anyhow!("Failed to manifest Jsonnet value: {}", e))?;

        serde_json::from_str(&json_str)
            .map_err(|e| anyhow!("Failed to parse manifested JSON: {}", e))
    }
}

/// A cache for reusable JsonnetEvaluators using thread-local storage
/// This allows reusing evaluators within the same thread without thread-safety issues
pub struct SharedEvaluatorPool {
    _marker: (),
}

thread_local! {
    static EVALUATOR_CACHE: RefCell<HashMap<String, JsonnetEvaluator>> = RefCell::new(HashMap::new());
}

impl SharedEvaluatorPool {
    /// Create a new shared evaluator pool
    pub fn new(_pool_size: usize) -> Self {
        use log::info;
        info!("Creating thread-local evaluator cache");
        SharedEvaluatorPool { _marker: () }
    }

    /// Use an evaluator with specific import paths, caching it for reuse
    pub fn with_evaluator<F, R>(&self, import_paths: &[PathBuf], f: F) -> Result<R>
    where
        F: FnOnce(&JsonnetEvaluator) -> Result<R>,
    {
        // Create a cache key from the import paths
        let cache_key = import_paths
            .iter()
            .map(|p| p.to_string_lossy())
            .collect::<Vec<_>>()
            .join(":");

        EVALUATOR_CACHE.with(|cache| {
            let mut cache = cache.borrow_mut();

            // Get or create an evaluator for these import paths
            if !cache.contains_key(&cache_key) {
                use log::debug;
                debug!("Creating new evaluator for import paths: {}", cache_key);
                let evaluator = JsonnetEvaluator::new_with_paths(import_paths.to_vec())?;
                cache.insert(cache_key.clone(), evaluator);
            }

            let evaluator = cache.get(&cache_key).unwrap();
            f(evaluator)
        })
    }
}

// Helper function to compute import paths for a file
fn compute_import_paths(path: &Path) -> Vec<PathBuf> {
    let mut import_paths = Vec::new();

    // Start from the file's directory and walk up to find lib/ and vendor/
    let start_dir = if path.is_dir() {
        path.to_path_buf()
    } else {
        path.parent()
            .map(|p| p.to_path_buf())
            .unwrap_or_else(|| PathBuf::from("."))
    };

    let mut current = start_dir.clone();
    loop {
        // Check for lib/ directory
        let lib_dir = current.join("lib");
        if lib_dir.is_dir() && !import_paths.contains(&lib_dir) {
            import_paths.push(lib_dir);
        }

        // Check for vendor/ directory
        let vendor_dir = current.join("vendor");
        if vendor_dir.is_dir() && !import_paths.contains(&vendor_dir) {
            import_paths.push(vendor_dir);
        }

        // Check for jsonnetfile.json (indicates project root)
        if current.join("jsonnetfile.json").exists()
            || current.join("jsonnetfile.lock.json").exists()
        {
            // Found project root, add it too
            if !import_paths.contains(&current) {
                import_paths.push(current.clone());
            }
            break;
        }

        // Move up one directory
        match current.parent() {
            Some(parent) => current = parent.to_path_buf(),
            None => break, // Reached filesystem root
        }
    }

    // Add the file's own directory
    if !import_paths.contains(&start_dir) {
        import_paths.push(start_dir);
    }

    import_paths
}

// Load all environments from a file with their full data
pub fn load_all_environments(path: &Path) -> Result<Vec<LoadedEnvironment>> {
    load_all_environments_with_pool(path, None)
}

// Load all environments with an optional evaluator pool for reuse
pub fn load_all_environments_with_pool(
    path: &Path,
    pool: Option<&SharedEvaluatorPool>,
) -> Result<Vec<LoadedEnvironment>> {
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

    // Evaluate with full data using pool if available
    let value = if let Some(pool) = pool {
        pool.with_evaluator(&import_paths, |evaluator| evaluator.eval_file(&main_file))?
    } else {
        let evaluator = JsonnetEvaluator::new_with_paths(import_paths)?;
        evaluator.eval_file(&main_file)?
    };

    // Try to parse as a single environment first
    if let Ok(environment) = serde_json::from_value::<Environment>(value.clone()) {
        return Ok(vec![LoadedEnvironment {
            path: main_file,
            environment,
        }]);
    }

    // Find all environments recursively in the structure
    let environments = find_environments_recursive(&value)?;

    if environments.is_empty() {
        return Err(anyhow!("No Environment objects found in file"));
    }

    // Return all environments found
    Ok(environments
        .into_iter()
        .map(|environment| LoadedEnvironment {
            path: main_file.clone(),
            environment,
        })
        .collect())
}

pub fn load_environment_by_name(path: &Path, name: &str) -> Result<LoadedEnvironment> {
    load_environment_by_name_with_pool(path, name, None)
}

// Load a specific environment by name with an optional evaluator pool
pub fn load_environment_by_name_with_pool(
    path: &Path,
    name: &str,
    pool: Option<&SharedEvaluatorPool>,
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

    // Use SingleEnvEvalScript to load only this environment with its data
    let script = SINGLE_ENV_EVAL_SCRIPT.replace("%s", name);
    let value = if let Some(pool) = pool {
        pool.with_evaluator(&import_paths, |evaluator| {
            evaluator.eval_script(&main_file, &script)
        })?
    } else {
        let evaluator = JsonnetEvaluator::new_with_paths(import_paths)?;
        evaluator.eval_script(&main_file, &script)?
    };

    // Find environments recursively in the structure
    let environments = find_environments_recursive(&value)?;

    if let Some(env) = environments.into_iter().next() {
        return Ok(LoadedEnvironment {
            path: main_file,
            environment: env,
        });
    }

    Err(anyhow!(
        "No environment named '{}' found in file {:?}",
        name,
        main_file
    ))
}

pub fn load_environment(path: &Path) -> Result<LoadedEnvironment> {
    load_environment_with_pool(path, None)
}

// Load an environment with an optional evaluator pool
pub fn load_environment_with_pool(
    path: &Path,
    pool: Option<&SharedEvaluatorPool>,
) -> Result<LoadedEnvironment> {
    // Check if this is a static environment (has spec.json)
    let env_dir = if path.is_dir() {
        path
    } else {
        path.parent().unwrap_or(path)
    };

    let spec_file = env_dir.join("spec.json");
    if spec_file.exists() {
        return load_static_environment_with_pool(env_dir, pool);
    }

    // Otherwise, load as inline environment
    load_inline_environment_with_pool(path, pool)
}

fn load_static_environment(path: &Path) -> Result<LoadedEnvironment> {
    load_static_environment_with_pool(path, None)
}

fn load_static_environment_with_pool(
    path: &Path,
    pool: Option<&SharedEvaluatorPool>,
) -> Result<LoadedEnvironment> {
    use log::debug;

    // Read spec.json
    let spec_file = path.join("spec.json");
    let spec_contents = fs::read_to_string(&spec_file)
        .with_context(|| format!("Failed to read spec.json: {:?}", spec_file))?;

    let mut environment: Environment = serde_json::from_str(&spec_contents)
        .with_context(|| format!("Failed to parse spec.json: {:?}", spec_file))?;

    debug!(
        "Loaded static environment from spec.json: {:?}",
        environment.metadata.name
    );

    // Populate metadata.namespace from spec.namespace if not set
    if environment.metadata.namespace.is_none() {
        environment.metadata.namespace = Some(environment.spec.namespace.clone());
    }

    let import_paths = compute_import_paths(path);

    // Look for main.jsonnet
    let main_file = path.join("main.jsonnet");
    if !main_file.exists() {
        return Err(anyhow!("No main.jsonnet found in directory: {:?}", path));
    }

    // Evaluate jsonnet to get data using pool if available
    let value = if let Some(pool) = pool {
        pool.with_evaluator(&import_paths, |evaluator| evaluator.eval_file(&main_file))?
    } else {
        let evaluator = JsonnetEvaluator::new_with_paths(import_paths)?;
        evaluator.eval_file(&main_file)?
    };
    environment.data = Some(value);

    Ok(LoadedEnvironment {
        path: main_file,
        environment,
    })
}

fn load_inline_environment(path: &Path) -> Result<LoadedEnvironment> {
    load_inline_environment_with_pool(path, None)
}

fn load_inline_environment_with_pool(
    path: &Path,
    pool: Option<&SharedEvaluatorPool>,
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

    let value = if let Some(pool) = pool {
        pool.with_evaluator(&import_paths, |evaluator| evaluator.eval_file(&main_file))?
    } else {
        let evaluator = JsonnetEvaluator::new_with_paths(import_paths)?;
        evaluator.eval_file(&main_file)?
    };

    // Try to parse as a single environment first
    if let Ok(environment) = serde_json::from_value::<Environment>(value.clone()) {
        return Ok(LoadedEnvironment {
            path: main_file,
            environment,
        });
    }

    // If that fails, try to find environments recursively in the structure
    // (e.g., in an "envs" field or other nested locations)
    let environments = find_environments_recursive(&value)?;

    if environments.is_empty() {
        return Err(anyhow!("No Environment objects found in file"));
    }

    if environments.len() > 1 {
        return Err(anyhow!(
            "File contains multiple environments. Use recursive mode to load them individually."
        ));
    }

    Ok(LoadedEnvironment {
        path: main_file,
        environment: environments.into_iter().next().unwrap(),
    })
}

// List all environments in a file (for recursive discovery)
pub fn list_environments(path: &Path) -> Result<Vec<LoadedEnvironment>> {
    list_environments_with_pool(path, None)
}

// List all environments with an optional evaluator pool
pub fn list_environments_with_pool(
    path: &Path,
    pool: Option<&SharedEvaluatorPool>,
) -> Result<Vec<LoadedEnvironment>> {
    // Check if this is a static environment (has spec.json)
    let env_dir = if path.is_dir() {
        path
    } else {
        path.parent().unwrap_or(path)
    };

    let spec_file = env_dir.join("spec.json");
    if spec_file.exists() {
        // For static environments, just load the metadata from spec.json
        let spec_contents = fs::read_to_string(&spec_file)
            .with_context(|| format!("Failed to read spec.json: {:?}", spec_file))?;

        let environment: Environment = serde_json::from_str(&spec_contents)
            .with_context(|| format!("Failed to parse spec.json: {:?}", spec_file))?;

        return Ok(vec![LoadedEnvironment {
            path: spec_file,
            environment,
        }]);
    }

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

    // Use metadata-only evaluation to avoid loading all manifest data
    // This is much faster as it doesn't evaluate the .data field
    let value = if let Some(pool) = pool {
        pool.with_evaluator(&import_paths, |evaluator| {
            evaluator.eval_script(&main_file, METADATA_EVAL_SCRIPT)
        })?
    } else {
        let evaluator = JsonnetEvaluator::new_with_paths(import_paths)?;
        evaluator.eval_script(&main_file, METADATA_EVAL_SCRIPT)?
    };

    // Try to parse as a single environment first
    if let Ok(environment) = serde_json::from_value::<Environment>(value.clone()) {
        return Ok(vec![LoadedEnvironment {
            path: main_file,
            environment,
        }]);
    }

    // If that fails, try to find environments recursively in the structure
    let environments = find_environments_recursive(&value)?;

    if environments.is_empty() {
        return Err(anyhow!("No Environment objects found in file"));
    }

    // Return all environments found
    Ok(environments
        .into_iter()
        .map(|environment| LoadedEnvironment {
            path: main_file.clone(),
            environment,
        })
        .collect())
}

pub fn load_manifests(env: &LoadedEnvironment) -> Result<Vec<Manifest>> {
    // Convert the environment's spec.data to JSON and extract manifests
    // The Environment object itself contains the data we need in its spec.data field
    let env_json =
        serde_json::to_value(&env.environment).context("Failed to serialize environment")?;

    // Extract manifests from the environment's data field
    let manifests = extract_manifests_from_value(&env_json)?;

    Ok(manifests)
}

fn extract_manifests_from_value(value: &Value) -> Result<Vec<Manifest>> {
    let mut manifests = Vec::new();

    // If the value has a "data" field, use that
    let data_value = if let Some(data) = value.get("data") {
        data
    } else {
        value
    };

    match data_value {
        Value::Object(obj) => {
            // Iterate through all keys in the object
            for (_key, val) in obj {
                if let Value::Object(manifest_obj) = val {
                    // Check if this looks like a Kubernetes manifest
                    if manifest_obj.contains_key("apiVersion") && manifest_obj.contains_key("kind")
                    {
                        let manifest = Manifest::new(
                            manifest_obj
                                .iter()
                                .map(|(k, v)| (k.clone(), v.clone()))
                                .collect(),
                        )?;
                        manifests.push(manifest);
                    } else {
                        // Recursively check nested objects
                        manifests.extend(extract_manifests_from_value(val)?);
                    }
                } else if val.is_object() || val.is_array() {
                    // Recursively check nested values
                    manifests.extend(extract_manifests_from_value(val)?);
                }
            }
        }
        Value::Array(arr) => {
            for item in arr {
                manifests.extend(extract_manifests_from_value(item)?);
            }
        }
        _ => {}
    }

    Ok(manifests)
}

// Recursively search for Environment objects in a JSON value
// This mimics Tanka's MetadataEvalScript behavior
fn find_environments_recursive(value: &Value) -> Result<Vec<Environment>> {
    let mut environments = Vec::new();

    match value {
        Value::Object(obj) => {
            // Check if this object is an Environment
            if let (Some(Value::String(kind)), Some(Value::String(_api_version))) =
                (obj.get("kind"), obj.get("apiVersion"))
            {
                if kind == "Environment" {
                    if let Ok(env) = serde_json::from_value::<Environment>(value.clone()) {
                        environments.push(env);
                        return Ok(environments);
                    }
                }
            }

            // Recursively search nested objects
            for (_key, nested_value) in obj {
                environments.extend(find_environments_recursive(nested_value)?);
            }
        }
        Value::Array(arr) => {
            // Recursively search array elements
            for item in arr {
                environments.extend(find_environments_recursive(item)?);
            }
        }
        _ => {}
    }

    Ok(environments)
}
