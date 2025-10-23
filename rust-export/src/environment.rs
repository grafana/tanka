use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::path::PathBuf;

#[derive(Debug, Clone, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Environment {
    pub api_version: String,
    pub kind: String,
    pub metadata: Metadata,
    pub spec: Spec,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub data: Option<serde_json::Value>,
}

impl Environment {
    pub fn new() -> Self {
        Environment {
            api_version: "tanka.dev/v1alpha1".to_string(),
            kind: "Environment".to_string(),
            metadata: Metadata::default(),
            spec: Spec::default(),
            data: None,
        }
    }
}

#[derive(Debug, Clone, Default, Serialize, Deserialize)]
pub struct Metadata {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub name: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub namespace: Option<String>,
    #[serde(default)]
    pub labels: HashMap<String, String>,
}

#[derive(Debug, Clone, Default, Serialize, Deserialize)]
#[serde(rename_all = "camelCase")]
pub struct Spec {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub api_server: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub context_names: Option<Vec<String>>,
    #[serde(default = "default_namespace")]
    pub namespace: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub diff_strategy: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub apply_strategy: Option<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub inject_labels: Option<bool>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub tanka_env_label_from_fields: Option<Vec<String>>,
    #[serde(default)]
    pub resource_defaults: ResourceDefaults,
    #[serde(default)]
    pub expect_versions: ExpectVersions,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub export_jsonnet_implementation: Option<String>,
}

fn default_namespace() -> String {
    "default".to_string()
}

#[derive(Debug, Clone, Default, Serialize, Deserialize)]
pub struct ResourceDefaults {
    #[serde(default)]
    pub annotations: HashMap<String, String>,
    #[serde(default)]
    pub labels: HashMap<String, String>,
}

#[derive(Debug, Clone, Default, Serialize, Deserialize)]
pub struct ExpectVersions {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub tanka: Option<String>,
}

/// Represents a Tanka environment with its file path
#[derive(Debug, Clone)]
pub struct LoadedEnvironment {
    pub path: PathBuf,
    pub environment: Environment,
}
