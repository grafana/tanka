use crate::loader::{detect_loader, LoaderOpts};
use serde_json::json;
use std::path::PathBuf;

#[test]
fn test_load_static() {
    let base_dir = PathBuf::from("testdata/cases/withspecjson");
    let loader = detect_loader(&base_dir, None).expect("Failed to detect loader");
    assert_eq!(loader.name(), "static");

    let loaded = loader
        .load(&base_dir, &LoaderOpts::default())
        .expect("Failed to load environment");

    // Check environment metadata
    assert_eq!(
        loaded.environment.api_version,
        "tanka.dev/v1alpha1",
        "API version mismatch"
    );
    assert_eq!(loaded.environment.kind, "Environment", "Kind mismatch");
    assert_eq!(
        loaded.environment.metadata.name,
        Some("withspec".to_string()),
        "Name mismatch"
    );
    assert_eq!(
        loaded.environment.spec.api_server,
        Some("https://localhost".to_string()),
        "API server mismatch"
    );
    assert_eq!(loaded.environment.spec.namespace, "withspec");

    // Check data (manifests)
    assert!(loaded.environment.data.is_some(), "Data is missing");
    let data = loaded.environment.data.unwrap();

    let expected_manifest = json!({
        "apiVersion": "v1",
        "kind": "ConfigMap",
        "metadata": {
            "name": "config"
        }
    });

    assert_eq!(data, expected_manifest, "Manifest mismatch");
}

#[test]
fn test_load_static_filename() {
    let base_dir = PathBuf::from("testdata/cases/withspecjson/main.jsonnet");
    let loader = detect_loader(&base_dir, None).expect("Failed to detect loader");
    assert_eq!(loader.name(), "static");

    let loaded = loader
        .load(&base_dir.parent().unwrap(), &LoaderOpts::default())
        .expect("Failed to load environment");

    // Check environment metadata
    assert_eq!(
        loaded.environment.api_version,
        "tanka.dev/v1alpha1",
        "API version mismatch"
    );
    assert_eq!(loaded.environment.kind, "Environment", "Kind mismatch");
    assert_eq!(
        loaded.environment.metadata.name,
        Some("withspec".to_string()),
        "Name mismatch"
    );

    // Check data
    assert!(loaded.environment.data.is_some(), "Data is missing");
}

#[test]
fn test_load_inline() {
    let base_dir = PathBuf::from("testdata/cases/withenv");
    let loader = detect_loader(&base_dir, None).expect("Failed to detect loader");
    assert_eq!(loader.name(), "inline");

    let loaded = loader
        .load(&base_dir, &LoaderOpts::default())
        .expect("Failed to load environment");

    // Check environment metadata
    assert_eq!(
        loaded.environment.api_version,
        "tanka.dev/v1alpha1",
        "API version mismatch"
    );
    assert_eq!(loaded.environment.kind, "Environment", "Kind mismatch");
    assert_eq!(
        loaded.environment.metadata.name,
        Some("withenv".to_string()),
        "Name mismatch"
    );
    assert_eq!(
        loaded.environment.spec.api_server,
        Some("https://localhost".to_string()),
        "API server mismatch"
    );
    assert_eq!(loaded.environment.spec.namespace, "withenv");

    // Check data (manifests)
    assert!(loaded.environment.data.is_some(), "Data is missing");
    let data = loaded.environment.data.unwrap();

    let expected_manifest = json!({
        "apiVersion": "v1",
        "kind": "ConfigMap",
        "metadata": {
            "name": "config"
        }
    });

    assert_eq!(data, expected_manifest, "Manifest mismatch");
}

#[test]
fn test_load_inline_filename() {
    let base_dir = PathBuf::from("testdata/cases/withenv/main.jsonnet");
    let loader = detect_loader(&base_dir, None).expect("Failed to detect loader");
    assert_eq!(loader.name(), "inline");

    let loaded = loader
        .load(&base_dir, &LoaderOpts::default())
        .expect("Failed to load environment");

    // Check environment metadata
    assert_eq!(
        loaded.environment.api_version,
        "tanka.dev/v1alpha1",
        "API version mismatch"
    );
    assert_eq!(loaded.environment.kind, "Environment", "Kind mismatch");
    assert_eq!(
        loaded.environment.metadata.name,
        Some("withenv".to_string()),
        "Name mismatch"
    );

    // Check data
    assert!(loaded.environment.data.is_some(), "Data is missing");
}

#[test]
fn test_peek_static() {
    let base_dir = PathBuf::from("testdata/cases/withspecjson");
    let loader = detect_loader(&base_dir, None).expect("Failed to detect loader");

    let env = loader
        .peek(&base_dir, &LoaderOpts::default())
        .expect("Failed to peek environment");

    // Peek should not load data
    assert!(env.data.is_none(), "Peek should not load data");
    assert_eq!(env.metadata.name, Some("withspec".to_string()));
}

#[test]
fn test_peek_inline() {
    let base_dir = PathBuf::from("testdata/cases/withenv");
    let loader = detect_loader(&base_dir, None).expect("Failed to detect loader");

    let env = loader
        .peek(&base_dir, &LoaderOpts::default())
        .expect("Failed to peek environment");

    // Peek should not load data
    assert!(env.data.is_none(), "Peek should not load data");
    assert_eq!(env.metadata.name, Some("withenv".to_string()));
}

#[test]
fn test_list_static() {
    let base_dir = PathBuf::from("testdata/cases/withspecjson");
    let loader = detect_loader(&base_dir, None).expect("Failed to detect loader");

    let envs = loader
        .list(&base_dir, &LoaderOpts::default())
        .expect("Failed to list environments");

    assert_eq!(envs.len(), 1, "Static should return exactly one environment");
    assert_eq!(envs[0].metadata.name, Some("withspec".to_string()));
}

#[test]
fn test_list_inline() {
    let base_dir = PathBuf::from("testdata/cases/withenv");
    let loader = detect_loader(&base_dir, None).expect("Failed to detect loader");

    let envs = loader
        .list(&base_dir, &LoaderOpts::default())
        .expect("Failed to list environments");

    assert_eq!(
        envs.len(),
        1,
        "Inline should return at least one environment"
    );
    assert_eq!(envs[0].metadata.name, Some("withenv".to_string()));
}

#[test]
fn test_eval_static() {
    let base_dir = PathBuf::from("testdata/cases/withspecjson");
    let loader = detect_loader(&base_dir, None).expect("Failed to detect loader");

    let data = loader
        .eval(&base_dir, &LoaderOpts::default())
        .expect("Failed to eval");

    // Should return raw jsonnet output (the ConfigMap)
    let expected = json!({
        "apiVersion": "v1",
        "kind": "ConfigMap",
        "metadata": {
            "name": "config"
        }
    });

    assert_eq!(data, expected);
}

#[test]
fn test_eval_inline() {
    let base_dir = PathBuf::from("testdata/cases/withenv");
    let loader = detect_loader(&base_dir, None).expect("Failed to detect loader");

    let data = loader
        .eval(&base_dir, &LoaderOpts::default())
        .expect("Failed to eval");

    // Should return the full Environment object
    assert!(data.is_object());
    let obj = data.as_object().unwrap();
    assert_eq!(obj.get("kind").unwrap(), "Environment");
}

