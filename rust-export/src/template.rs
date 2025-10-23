use anyhow::{anyhow, Context, Result};
use gtmpl::{template, Context as GtmplContext, Func, Template, Value as GtmplValue};
use serde_json::Value as JsonValue;

use crate::environment::Environment;
use crate::manifest::Manifest;

const BEL_RUNE: &str = "\x07";

pub struct TemplateEngine {
    template_str: String,
}

impl TemplateEngine {
    pub fn new(format: &str, _environment: &Environment) -> Result<Self> {
        // Replace path separators with BEL character in the template
        let replaced_format =
            replace_in_template_text(format, std::path::MAIN_SEPARATOR_STR, BEL_RUNE);

        // Note: We don't validate the template here because it may reference variables
        // that won't be available until rendering time (like 'env'). Errors will be
        // caught during the actual rendering.

        Ok(TemplateEngine {
            template_str: replaced_format,
        })
    }

    pub fn apply(&self, manifest: &Manifest, environment: &Environment) -> Result<String> {
        use gtmpl::{Context as GtmplContext, Template};
        use std::collections::HashMap;

        // Convert manifest and environment to gtmpl values
        let mut manifest_json =
            serde_json::to_value(manifest).context("Failed to serialize manifest")?;
        let mut env_json =
            serde_json::to_value(environment).context("Failed to serialize environment")?;

        // gtmpl is strict about missing fields, unlike Go templates which return nil/false
        // Add default empty values for commonly used fields to prevent errors
        let common_label_fields = vec!["fluxExport", "fluxExportDir"];

        // Add defaults to environment labels
        if let Some(env_obj) = env_json.as_object_mut() {
            if let Some(metadata) = env_obj.get_mut("metadata").and_then(|m| m.as_object_mut()) {
                if let Some(labels) = metadata.get_mut("labels").and_then(|l| l.as_object_mut()) {
                    for field in &common_label_fields {
                        labels
                            .entry(field.to_string())
                            .or_insert(JsonValue::String("".to_string()));
                    }
                }
            }
        }

        // Add defaults to manifest labels
        if let Some(manifest_obj) = manifest_json.as_object_mut() {
            if let Some(metadata) = manifest_obj
                .get_mut("metadata")
                .and_then(|m| m.as_object_mut())
            {
                // Ensure labels exists
                if !metadata.contains_key("labels") {
                    metadata.insert(
                        "labels".to_string(),
                        JsonValue::Object(serde_json::Map::new()),
                    );
                }
                if let Some(labels) = metadata.get_mut("labels").and_then(|l| l.as_object_mut()) {
                    for field in &common_label_fields {
                        labels
                            .entry(field.to_string())
                            .or_insert(JsonValue::String("".to_string()));
                    }
                }
            }
        }

        // Combine manifest and environment into a single context
        // The manifest fields are at root level, and env is available as a field
        let mut combined_data = if let JsonValue::Object(map) = manifest_json {
            map
        } else {
            serde_json::Map::new()
        };
        combined_data.insert("env".to_string(), env_json);

        let combined_value = json_to_gtmpl_value(&JsonValue::Object(combined_data));

        // Create template (Go template-compatible)
        // Note: In Go templates, 'env' is registered as a function and can be called without a dot.
        // In gtmpl (Rust), we need to use '.env' to access it as a field.
        // Rewrite the template to convert 'env.' to '.env.' for compatibility.
        // We need to handle cases like: {{ env.field }}, {{ if env.field }}, {{ if not env.field }}, etc.
        use regex::Regex;
        let env_regex = Regex::new(r"\benv\.").unwrap();
        let compat_template = env_regex
            .replace_all(&self.template_str, ".env.")
            .into_owned();

        let mut tmpl = Template::default();
        // Enable option:missingkey=zero to handle missing fields like Go templates do
        // This makes missing fields return their zero value instead of errors
        tmpl.parse("{{- /*gotype: .*/}}\n").ok(); // Dummy parse to initialize
        tmpl.parse(&compat_template)
            .map_err(|e| anyhow!("Failed to parse template: {}", e))?;

        // Create context from combined value
        let mut context = GtmplContext::from(combined_value)
            .map_err(|e| anyhow!("Failed to create context: {}", e))?;

        // Try to render with lenient error handling
        let result = match tmpl.render(&context) {
            Ok(r) => r,
            Err(e) if e.contains("no field") => {
                // If we get a missing field error, this might be expected in the template logic
                // For now, return a descriptive error. In the future, we could add default values.
                return Err(anyhow!("Template requires a field that doesn't exist: {}. \
                                   Consider updating your template to check for field existence first, \
                                   or add the field to your environment metadata.", e));
            }
            Err(e) => return Err(anyhow!("Failed to render filename template: {}", e)),
        };

        // Replace path separators with dashes to avoid accidental subdirectories
        let result = result.replace(std::path::MAIN_SEPARATOR, "-");

        // Replace BEL character back with path separator for intentional subdirectories
        let result = result.replace(BEL_RUNE, std::path::MAIN_SEPARATOR_STR);

        Ok(result)
    }

    pub fn render_filename(
        &self,
        env: &crate::environment::LoadedEnvironment,
        manifest: &Manifest,
    ) -> Result<String> {
        self.apply(manifest, &env.environment)
    }
}

/// Convert serde_json::Value to gtmpl::Value
fn json_to_gtmpl_value(json: &JsonValue) -> GtmplValue {
    match json {
        JsonValue::Null => GtmplValue::Nil,
        JsonValue::Bool(b) => GtmplValue::Bool(*b),
        JsonValue::Number(n) => {
            // Convert to f64 first, then to Number
            let f = n.as_f64().unwrap_or(0.0);
            GtmplValue::Number(f.into())
        }
        JsonValue::String(s) => GtmplValue::String(s.clone()),
        JsonValue::Array(arr) => {
            let vals: Vec<GtmplValue> = arr.iter().map(json_to_gtmpl_value).collect();
            GtmplValue::Array(vals)
        }
        JsonValue::Object(obj) => {
            let mut map = std::collections::HashMap::new();
            for (k, v) in obj {
                map.insert(k.clone(), json_to_gtmpl_value(v));
            }
            GtmplValue::Object(map)
        }
    }
}

/// Replace text within template delimiters
fn replace_in_template_text(template_str: &str, from: &str, to: &str) -> String {
    let mut result = String::new();
    let mut in_template = false;
    let mut chars = template_str.chars().peekable();

    while let Some(c) = chars.next() {
        if c == '{' && chars.peek() == Some(&'{') {
            // Entering template
            result.push(c);
            result.push(chars.next().unwrap());
            in_template = true;
        } else if c == '}' && chars.peek() == Some(&'}') {
            // Exiting template
            result.push(c);
            result.push(chars.next().unwrap());
            in_template = false;
        } else if in_template && c.to_string() == from {
            // Replace separator inside template
            result.push_str(to);
        } else {
            result.push(c);
        }
    }

    result
}
