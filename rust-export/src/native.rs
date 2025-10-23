use anyhow::{anyhow, Context, Result};
use jrsonnet_evaluator::{
    function::{builtin, CallLocation},
    IStr, ObjValue, Result as JResult, Val,
};
use jrsonnet_stdlib::ContextInitializer;
use log::debug;
use serde::{Deserialize, Serialize};
use serde_json::Value as JsonValue;
use std::collections::HashMap;
use std::io::Write;
use std::process::{Command, Stdio};

/// Register all Tanka native functions with the ContextInitializer
pub fn register_native_functions(ctx_init: &ContextInitializer) {
    // Create and register all native functions
    ctx_init
        .settings_mut()
        .ext_natives
        .insert("helmTemplate".into(), builtin_helm_template::INST.into());

    ctx_init
        .settings_mut()
        .ext_natives
        .insert("parseJson".into(), builtin_parse_json::INST.into());

    ctx_init
        .settings_mut()
        .ext_natives
        .insert("parseYaml".into(), builtin_parse_yaml::INST.into());

    ctx_init.settings_mut().ext_natives.insert(
        "manifestJsonFromJson".into(),
        builtin_manifest_json_from_json::INST.into(),
    );

    ctx_init.settings_mut().ext_natives.insert(
        "manifestYamlFromJson".into(),
        builtin_manifest_yaml_from_json::INST.into(),
    );

    ctx_init.settings_mut().ext_natives.insert(
        "escapeStringRegex".into(),
        builtin_escape_string_regex::INST.into(),
    );

    ctx_init
        .settings_mut()
        .ext_natives
        .insert("regexMatch".into(), builtin_regex_match::INST.into());

    ctx_init
        .settings_mut()
        .ext_natives
        .insert("regexSubst".into(), builtin_regex_subst::INST.into());

    ctx_init
        .settings_mut()
        .ext_natives
        .insert("sha256".into(), builtin_sha256::INST.into());
}

#[derive(Debug, Deserialize)]
#[serde(rename_all = "camelCase")]
struct HelmOpts {
    called_from: Option<String>,
    values: Option<HashMap<String, JsonValue>>,
    api_versions: Option<Vec<String>>,
    #[serde(default = "default_true")]
    include_crds: bool,
    skip_tests: Option<bool>,
    kube_version: Option<String>,
    namespace: Option<String>,
    no_hooks: Option<bool>,
}

fn default_true() -> bool {
    true
}

/// Implements the helmTemplate native function
/// This calls the helm CLI to template Kubernetes manifests
#[builtin]
pub fn builtin_helm_template(
    _loc: CallLocation,
    name: IStr,
    chart: IStr,
    opts: ObjValue,
) -> JResult<Val> {
    use jrsonnet_evaluator::{error::ErrorKind, ObjValueBuilder};

    let name_str = name.to_string();
    let chart_str = chart.to_string();

    // Convert ObjValue to JSON to parse into HelmOpts
    let opts_json = obj_to_json(&opts)?;
    let helm_opts: HelmOpts = serde_json::from_value(opts_json).map_err(|e| {
        ErrorKind::RuntimeError(format!("Failed to parse helm options: {}", e).into())
    })?;

    if helm_opts.called_from.is_none() {
        return Err(ErrorKind::RuntimeError("helmTemplate: 'opts.calledFrom' is unset or empty.\nTanka needs this to find your charts. See https://tanka.dev/helm#optscalledfrom-unset".into()).into());
    }

    // Execute helm template
    let manifests = execute_helm_template(&name_str, &chart_str, &helm_opts)
        .map_err(|e| ErrorKind::RuntimeError(format!("helmTemplate failed: {}", e).into()))?;

    debug!("Helm template returned {} manifests", manifests.len());

    // Convert manifests to a map - return as an object with keys based on kind and name
    let mut result = ObjValueBuilder::new();
    for (idx, manifest) in manifests.iter().enumerate() {
        if let Some(obj) = manifest.as_object() {
            let kind = obj
                .get("kind")
                .and_then(|v| v.as_str())
                .unwrap_or("unknown");
            let name = obj
                .get("metadata")
                .and_then(|m| m.as_object())
                .and_then(|m| m.get("name"))
                .and_then(|n| n.as_str())
                .unwrap_or("unknown");

            // Convert name format: replace hyphens with underscores to match Tanka conventions
            // e.g. "kustomize-controller" -> "kustomize_controller"
            let normalized_name = name.replace('-', "_");
            let key = format!("{}_{}", kind.to_lowercase(), normalized_name);

            // Convert manifest to Val by re-parsing through JSON
            let manifest_str = serde_json::to_string(manifest).map_err(|e| {
                ErrorKind::RuntimeError(format!("Failed to serialize manifest: {}", e).into())
            })?;
            let manifest_val: serde_json::Value =
                serde_json::from_str(&manifest_str).map_err(|e| {
                    ErrorKind::RuntimeError(format!("Failed to parse manifest: {}", e).into())
                })?;

            let val = json_to_jsonnet_val(&manifest_val)?;
            result.field(&key).value(val);
        } else {
            let key = format!("manifest_{}", idx);
            let manifest_str = serde_json::to_string(manifest).map_err(|e| {
                ErrorKind::RuntimeError(format!("Failed to serialize manifest: {}", e).into())
            })?;
            let manifest_val: serde_json::Value =
                serde_json::from_str(&manifest_str).map_err(|e| {
                    ErrorKind::RuntimeError(format!("Failed to parse manifest: {}", e).into())
                })?;

            let val = json_to_jsonnet_val(&manifest_val)?;
            result.field(&key).value(val);
        }
    }

    Ok(Val::Obj(result.build()))
}

fn execute_helm_template(name: &str, chart: &str, opts: &HelmOpts) -> Result<Vec<JsonValue>> {
    use std::path::Path;

    debug!("Executing helm template: name={}, chart={}", name, chart);

    // Determine the working directory based on calledFrom
    // calledFrom is the file that called helmTemplate
    let working_dir = if let Some(called_from) = &opts.called_from {
        let called_from_path = Path::new(called_from);
        if let Some(parent) = called_from_path.parent() {
            parent.to_path_buf()
        } else {
            std::env::current_dir().context("Failed to get current directory")?
        }
    } else {
        std::env::current_dir().context("Failed to get current directory")?
    };

    debug!("Helm working directory: {:?}", working_dir);

    // Resolve chart path relative to the caller's directory
    // Strip leading slash if present to ensure relative path behavior
    let chart_clean = chart.strip_prefix('/').unwrap_or(chart);
    let chart_path = working_dir.join(chart_clean);

    // Verify the chart exists
    if !chart_path.exists() {
        return Err(anyhow!(
            "helmTemplate: Failed to find a chart at '{}'. See https://tanka.dev/helm#failed-to-find-chart",
            chart_path.display()
        ));
    }

    debug!("Resolved chart path: {:?}", chart_path);

    let mut args = vec![
        "template".to_string(),
        name.to_string(),
        chart_path.to_string_lossy().to_string(),
        "--values".to_string(),
        "-".to_string(), // Read values from stdin
    ];

    // Add flags
    if let Some(api_versions) = &opts.api_versions {
        for version in api_versions {
            args.push(format!("--api-versions={}", version));
        }
    }

    if opts.include_crds {
        args.push("--include-crds".to_string());
    }

    if opts.skip_tests.unwrap_or(false) {
        args.push("--skip-tests".to_string());
    }

    if let Some(kube_version) = &opts.kube_version {
        args.push(format!("--kube-version={}", kube_version));
    }

    if opts.no_hooks.unwrap_or(false) {
        args.push("--no-hooks".to_string());
    }

    if let Some(namespace) = &opts.namespace {
        args.push(format!("--namespace={}", namespace));
    }

    debug!("Helm command: helm {}", args.join(" "));

    let mut cmd = Command::new("helm")
        .args(&args)
        .stdin(Stdio::piped())
        .stdout(Stdio::piped())
        .stderr(Stdio::piped())
        .spawn()
        .context("Failed to spawn helm command")?;

    // Write values to stdin
    if let Some(stdin) = cmd.stdin.as_mut() {
        let values = opts.values.as_ref().map(|v| v.clone()).unwrap_or_default();
        let values_yaml =
            serde_yaml::to_string(&values).context("Failed to serialize values to YAML")?;
        stdin
            .write_all(values_yaml.as_bytes())
            .context("Failed to write values to helm stdin")?;
    }

    let output = cmd
        .wait_with_output()
        .context("Failed to wait for helm command")?;

    if !output.status.success() {
        let stderr = String::from_utf8_lossy(&output.stderr);
        return Err(anyhow!("helm template failed: {}", stderr));
    }

    let stdout = String::from_utf8(output.stdout).context("helm output is not valid UTF-8")?;

    debug!(
        "Helm output ({} bytes):\n{}",
        stdout.len(),
        if stdout.len() > 500 {
            &stdout[..500]
        } else {
            &stdout
        }
    );

    // Parse YAML documents
    let mut manifests = Vec::new();
    let deserializer = serde_yaml::Deserializer::from_str(&stdout);

    for document in deserializer {
        let value: JsonValue =
            JsonValue::deserialize(document).context("Failed to parse YAML document")?;

        // Skip empty documents
        if value.is_null() || (value.is_object() && value.as_object().unwrap().is_empty()) {
            continue;
        }

        debug!(
            "Parsed helm manifest: kind={:?}, name={:?}",
            value.get("kind"),
            value.get("metadata").and_then(|m| m.get("name"))
        );

        manifests.push(value);
    }

    debug!(
        "Helm template returned {} non-empty manifests",
        manifests.len()
    );

    Ok(manifests)
}

fn obj_to_json(obj: &ObjValue) -> JResult<JsonValue> {
    // Convert ObjValue to JSON
    // This is a simplified implementation
    let mut map = serde_json::Map::new();

    for (key, _) in obj.iter() {
        if let Some(value) = obj.get(key.clone())? {
            let json_val = val_to_json(&value)?;
            map.insert(key.to_string(), json_val);
        }
    }

    Ok(JsonValue::Object(map))
}

fn val_to_json(val: &Val) -> JResult<JsonValue> {
    match val {
        Val::Null => Ok(JsonValue::Null),
        Val::Bool(b) => Ok(JsonValue::Bool(*b)),
        Val::Num(n) => {
            // NumValue wraps an f64. We need to convert it via formatting and parsing
            // as there's no public getter
            let num_str = format!("{}", n);
            let num_f64 = num_str.parse::<f64>().unwrap_or(0.0);
            Ok(JsonValue::Number(
                serde_json::Number::from_f64(num_f64).unwrap_or(serde_json::Number::from(0)),
            ))
        }
        Val::Str(s) => Ok(JsonValue::String(s.to_string())),
        Val::Arr(arr) => {
            let mut result = Vec::new();
            for item in arr.iter_lazy() {
                let evaluated = item.evaluate()?;
                result.push(val_to_json(&evaluated)?);
            }
            Ok(JsonValue::Array(result))
        }
        Val::Obj(obj) => obj_to_json(obj),
        _ => Ok(JsonValue::String(format!("{:?}", val))),
    }
}

// Convert JSON value to jrsonnet Val using the stdlib parse function
fn json_to_jsonnet_val(json: &JsonValue) -> JResult<Val> {
    use jrsonnet_evaluator::error::ErrorKind;
    use jrsonnet_stdlib::builtin_parse_json;

    // Convert to JSON string and parse it using jrsonnet's builtin
    let json_str = serde_json::to_string(json)
        .map_err(|e| ErrorKind::RuntimeError(format!("Failed to serialize JSON: {}", e).into()))?;

    builtin_parse_json(json_str.into())
}

/// parseJson native function
#[builtin]
pub fn builtin_parse_json(_loc: CallLocation, json: IStr) -> JResult<Val> {
    use jrsonnet_stdlib::builtin_parse_json as parse;
    parse(json)
}

/// parseYaml native function
#[builtin]
pub fn builtin_parse_yaml(_loc: CallLocation, yaml: IStr) -> JResult<Val> {
    use jrsonnet_evaluator::error::ErrorKind;

    let yaml_str = yaml.to_string();
    let deserializer = serde_yaml::Deserializer::from_str(&yaml_str);
    let mut results = Vec::new();

    for document in deserializer {
        let value: JsonValue = JsonValue::deserialize(document)
            .map_err(|e| ErrorKind::RuntimeError(format!("Failed to parse YAML: {}", e).into()))?;
        results.push(value);
    }

    // Convert Vec<JsonValue> to Val
    let json_array = JsonValue::Array(results);
    json_to_jsonnet_val(&json_array)
}

/// manifestJsonFromJson native function
#[builtin]
pub fn builtin_manifest_json_from_json(
    _loc: CallLocation,
    json: IStr,
    indent: f64,
) -> JResult<Val> {
    use jrsonnet_evaluator::error::ErrorKind;

    let json_str = json.to_string();
    let indent_size = indent as usize;

    // Parse and re-serialize with proper indentation
    let value: JsonValue = serde_json::from_str(&json_str)
        .map_err(|e| ErrorKind::RuntimeError(format!("Failed to parse JSON: {}", e).into()))?;

    let mut output = Vec::new();
    let indent_bytes = vec![b' '; indent_size];
    let formatter = serde_json::ser::PrettyFormatter::with_indent(indent_bytes.as_slice());
    let mut serializer = serde_json::Serializer::with_formatter(&mut output, formatter);
    value
        .serialize(&mut serializer)
        .map_err(|e| ErrorKind::RuntimeError(format!("Failed to serialize JSON: {}", e).into()))?;

    output.push(b'\n');
    let result = String::from_utf8(output).map_err(|e| {
        ErrorKind::RuntimeError(format!("Failed to convert to UTF-8: {}", e).into())
    })?;

    Ok(Val::Str(result.into()))
}

/// manifestYamlFromJson native function
#[builtin]
pub fn builtin_manifest_yaml_from_json(_loc: CallLocation, json: IStr) -> JResult<Val> {
    use jrsonnet_evaluator::error::ErrorKind;

    let json_str = json.to_string();

    // Parse JSON
    let value: JsonValue = serde_json::from_str(&json_str)
        .map_err(|e| ErrorKind::RuntimeError(format!("Failed to parse JSON: {}", e).into()))?;

    // Serialize as YAML
    let yaml_str = serde_yaml::to_string(&value)
        .map_err(|e| ErrorKind::RuntimeError(format!("Failed to serialize YAML: {}", e).into()))?;

    Ok(Val::Str(yaml_str.into()))
}

/// escapeStringRegex native function
#[builtin]
pub fn builtin_escape_string_regex(_loc: CallLocation, s: IStr) -> JResult<Val> {
    use regex::escape;
    let escaped = escape(&s.to_string());
    Ok(Val::Str(escaped.into()))
}

/// regexMatch native function
#[builtin]
pub fn builtin_regex_match(_loc: CallLocation, regex: IStr, string: IStr) -> JResult<Val> {
    use jrsonnet_evaluator::error::ErrorKind;
    use regex::Regex;

    let re = Regex::new(&regex.to_string())
        .map_err(|e| ErrorKind::RuntimeError(format!("Invalid regex: {}", e).into()))?;

    Ok(Val::Bool(re.is_match(&string.to_string())))
}

/// regexSubst native function
#[builtin]
pub fn builtin_regex_subst(_loc: CallLocation, regex: IStr, src: IStr, repl: IStr) -> JResult<Val> {
    use jrsonnet_evaluator::error::ErrorKind;
    use regex::Regex;

    let re = Regex::new(&regex.to_string())
        .map_err(|e| ErrorKind::RuntimeError(format!("Invalid regex: {}", e).into()))?;

    let result = re
        .replace_all(&src.to_string(), repl.to_string().as_str())
        .to_string();
    Ok(Val::Str(result.into()))
}

/// sha256 native function
#[builtin]
pub fn builtin_sha256(_loc: CallLocation, s: IStr) -> JResult<Val> {
    use sha2::{Digest, Sha256};

    let mut hasher = Sha256::new();
    hasher.update(s.to_string().as_bytes());
    let result = hasher.finalize();
    let hex_string = format!("{:x}", result);

    Ok(Val::Str(hex_string.into()))
}
