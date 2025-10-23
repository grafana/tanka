use anyhow::{anyhow, Result};
use serde::{Deserialize, Serialize};
use serde_json::Value;
use std::collections::HashMap;

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Manifest {
    #[serde(flatten)]
    pub data: HashMap<String, Value>,
}

impl Manifest {
    pub fn new(data: HashMap<String, Value>) -> Result<Self> {
        let manifest = Manifest { data };
        manifest.verify()?;
        Ok(manifest)
    }

    pub fn verify(&self) -> Result<()> {
        // Check required fields
        if !self.data.contains_key("kind") {
            return Err(anyhow!("manifest missing 'kind' field"));
        }
        if !self.data.contains_key("apiVersion") {
            return Err(anyhow!("manifest missing 'apiVersion' field"));
        }

        // Check if it's a list type
        if !self.is_list() {
            // Non-list objects need metadata
            if !self.data.contains_key("metadata") {
                return Err(anyhow!("manifest missing 'metadata' field"));
            }

            // Check for name or generateName
            if let Some(metadata) = self.data.get("metadata") {
                if let Some(meta_obj) = metadata.as_object() {
                    if !meta_obj.contains_key("name") && !meta_obj.contains_key("generateName") {
                        return Err(anyhow!(
                            "manifest missing 'metadata.name' or 'metadata.generateName'"
                        ));
                    }
                }
            }
        }

        Ok(())
    }

    pub fn is_list(&self) -> bool {
        self.data
            .get("items")
            .map(|v| v.is_array())
            .unwrap_or(false)
    }

    pub fn kind(&self) -> Option<&str> {
        self.data.get("kind")?.as_str()
    }

    pub fn api_version(&self) -> Option<&str> {
        self.data.get("apiVersion")?.as_str()
    }

    pub fn metadata(&self) -> ManifestMetadata {
        ManifestMetadata::new(self.data.get("metadata"))
    }

    pub fn to_yaml(&self) -> Result<String> {
        Ok(serde_yaml::to_string(&self.data)?)
    }
}

#[derive(Debug)]
pub struct ManifestMetadata<'a> {
    data: Option<&'a Value>,
}

impl<'a> ManifestMetadata<'a> {
    fn new(data: Option<&'a Value>) -> Self {
        ManifestMetadata { data }
    }

    pub fn name(&self) -> Option<&str> {
        self.data?.as_object()?.get("name")?.as_str()
    }

    pub fn generate_name(&self) -> Option<&str> {
        self.data?.as_object()?.get("generateName")?.as_str()
    }

    pub fn namespace(&self) -> Option<&str> {
        self.data?.as_object()?.get("namespace")?.as_str()
    }

    pub fn labels(&self) -> HashMap<String, String> {
        let mut result = HashMap::new();
        if let Some(labels) = self
            .data
            .and_then(|v| v.as_object())
            .and_then(|obj| obj.get("labels"))
            .and_then(|v| v.as_object())
        {
            for (k, v) in labels {
                if let Some(s) = v.as_str() {
                    result.insert(k.clone(), s.to_string());
                }
            }
        }
        result
    }

    pub fn annotations(&self) -> HashMap<String, String> {
        let mut result = HashMap::new();
        if let Some(annotations) = self
            .data
            .and_then(|v| v.as_object())
            .and_then(|obj| obj.get("annotations"))
            .and_then(|v| v.as_object())
        {
            for (k, v) in annotations {
                if let Some(s) = v.as_str() {
                    result.insert(k.clone(), s.to_string());
                }
            }
        }
        result
    }
}
