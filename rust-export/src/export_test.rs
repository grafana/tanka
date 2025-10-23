use std::collections::HashMap;
use std::fs;
use std::path::{Path, PathBuf};

use anyhow::Result;
use walkdir::WalkDir;

use crate::export::{export_environments, ExportMergeStrategy, ExportOptions};

fn check_files(dir: &Path, expected_files: &[&str]) -> Result<()> {
    let mut existing_files = Vec::new();
    for entry in WalkDir::new(dir).into_iter().filter_map(|e| e.ok()) {
        if entry.file_type().is_file() {
            existing_files.push(entry.path().to_path_buf());
        }
    }

    let expected: Vec<PathBuf> = expected_files.iter().map(PathBuf::from).collect();

    // Sort for comparison
    let mut existing_sorted = existing_files.clone();
    existing_sorted.sort();
    let mut expected_sorted = expected.clone();
    expected_sorted.sort();

    if existing_sorted != expected_sorted {
        eprintln!("Expected files:");
        for f in &expected_sorted {
            eprintln!("  {:?}", f);
        }
        eprintln!("Actual files:");
        for f in &existing_sorted {
            eprintln!("  {:?}", f);
        }
        panic!("File lists don't match");
    }

    Ok(())
}

#[test]
fn test_export_environments() -> Result<()> {
    let temp_dir = tempfile::tempdir()?;
    let output_dir = temp_dir.path().to_path_buf();

    // Change to testdata directory
    let original_dir = std::env::current_dir()?;
    let testdata_dir = original_dir.join("testdata").canonicalize()?;
    std::env::set_current_dir(&testdata_dir)?;

    let result = (|| -> Result<()> {
        // Export all envs
        let paths = vec![PathBuf::from("test-export-envs")];
        let opts = ExportOptions {
            format: "{{.metadata.namespace}}/{{.metadata.name}}".to_string(),
            extension: "yaml".to_string(),
            parallelism: 1,
            cache_path: None,
            cache_path_regexes: vec![],
            merge_strategy: ExportMergeStrategy::None,
            merge_deleted_envs: vec![],
            skip_manifest: false,
            targets: vec![],
            name: None,
            selector: None,
        };

        export_environments(output_dir.clone(), paths, true, opts)?;

        // Check exported files
        check_files(
            &output_dir,
            &[
                &format!(
                    "{}/inline-namespace1/my-configmap.yaml",
                    output_dir.display()
                ),
                &format!(
                    "{}/inline-namespace1/my-deployment.yaml",
                    output_dir.display()
                ),
                &format!("{}/inline-namespace1/my-service.yaml", output_dir.display()),
                &format!(
                    "{}/inline-namespace2/my-deployment.yaml",
                    output_dir.display()
                ),
                &format!("{}/inline-namespace2/my-service.yaml", output_dir.display()),
                &format!("{}/static/initial-deployment.yaml", output_dir.display()),
                &format!("{}/static/initial-service.yaml", output_dir.display()),
                &format!("{}/manifest.json", output_dir.display()),
            ],
        )?;

        // Check manifest.json content
        let manifest_path = output_dir.join("manifest.json");
        let manifest_content = fs::read_to_string(&manifest_path)?;
        let manifest: HashMap<String, String> = serde_json::from_str(&manifest_content)?;

        let mut expected_manifest = HashMap::new();
        expected_manifest.insert(
            "inline-namespace1/my-configmap.yaml".to_string(),
            "test-export-envs/inline-envs/main.jsonnet".to_string(),
        );
        expected_manifest.insert(
            "inline-namespace1/my-deployment.yaml".to_string(),
            "test-export-envs/inline-envs/main.jsonnet".to_string(),
        );
        expected_manifest.insert(
            "inline-namespace1/my-service.yaml".to_string(),
            "test-export-envs/inline-envs/main.jsonnet".to_string(),
        );
        expected_manifest.insert(
            "inline-namespace2/my-deployment.yaml".to_string(),
            "test-export-envs/inline-envs/main.jsonnet".to_string(),
        );
        expected_manifest.insert(
            "inline-namespace2/my-service.yaml".to_string(),
            "test-export-envs/inline-envs/main.jsonnet".to_string(),
        );
        expected_manifest.insert(
            "static/initial-deployment.yaml".to_string(),
            "test-export-envs/static-env/main.jsonnet".to_string(),
        );
        expected_manifest.insert(
            "static/initial-service.yaml".to_string(),
            "test-export-envs/static-env/main.jsonnet".to_string(),
        );

        assert_eq!(
            manifest, expected_manifest,
            "Manifest content doesn't match"
        );

        // Try to re-export - should fail
        let paths = vec![PathBuf::from("test-export-envs")];
        let opts = ExportOptions {
            format: "{{.metadata.namespace}}/{{.metadata.name}}".to_string(),
            extension: "yaml".to_string(),
            parallelism: 1,
            cache_path: None,
            cache_path_regexes: vec![],
            merge_strategy: ExportMergeStrategy::None,
            merge_deleted_envs: vec![],
            skip_manifest: false,
            targets: vec![],
            name: None,
            selector: None,
        };

        let result = export_environments(output_dir.clone(), paths, true, opts);
        assert!(result.is_err(), "Should fail when directory is not empty");
        assert!(
            result.unwrap_err().to_string().contains("not empty"),
            "Error should mention directory not empty"
        );

        // Try with fail-on-conflicts
        let paths = vec![PathBuf::from("test-export-envs")];
        let opts = ExportOptions {
            format: "{{.metadata.namespace}}/{{.metadata.name}}".to_string(),
            extension: "yaml".to_string(),
            parallelism: 1,
            cache_path: None,
            cache_path_regexes: vec![],
            merge_strategy: ExportMergeStrategy::FailOnConflicts,
            merge_deleted_envs: vec![],
            skip_manifest: false,
            targets: vec![],
            name: None,
            selector: None,
        };

        let result = export_environments(output_dir.clone(), paths, true, opts);
        assert!(result.is_err(), "Should fail on conflicts");
        assert!(
            result.unwrap_err().to_string().contains("already exists"),
            "Error should mention file already exists"
        );

        // Re-export only static env with replace-envs strategy
        let paths = vec![PathBuf::from("test-export-envs/static-env")];
        let opts = ExportOptions {
            format: "{{.metadata.namespace}}/{{.metadata.name}}".to_string(),
            extension: "yaml".to_string(),
            parallelism: 1,
            cache_path: None,
            cache_path_regexes: vec![],
            merge_strategy: ExportMergeStrategy::ReplaceEnvs,
            merge_deleted_envs: vec![],
            skip_manifest: false,
            targets: vec![],
            name: None,
            selector: None,
        };

        export_environments(output_dir.clone(), paths, false, opts)?;

        check_files(
            &output_dir,
            &[
                &format!(
                    "{}/inline-namespace1/my-configmap.yaml",
                    output_dir.display()
                ),
                &format!(
                    "{}/inline-namespace1/my-deployment.yaml",
                    output_dir.display()
                ),
                &format!("{}/inline-namespace1/my-service.yaml", output_dir.display()),
                &format!(
                    "{}/inline-namespace2/my-deployment.yaml",
                    output_dir.display()
                ),
                &format!("{}/inline-namespace2/my-service.yaml", output_dir.display()),
                &format!("{}/static/initial-deployment.yaml", output_dir.display()),
                &format!("{}/static/initial-service.yaml", output_dir.display()),
                &format!("{}/manifest.json", output_dir.display()),
            ],
        )?;

        // Re-export and delete the files of inline envs
        let paths = vec![PathBuf::from("test-export-envs/static-env")];
        let opts = ExportOptions {
            format: "{{.metadata.namespace}}/{{.metadata.name}}".to_string(),
            extension: "yaml".to_string(),
            parallelism: 1,
            cache_path: None,
            cache_path_regexes: vec![],
            merge_strategy: ExportMergeStrategy::ReplaceEnvs,
            merge_deleted_envs: vec!["test-export-envs/inline-envs/main.jsonnet".to_string()],
            skip_manifest: false,
            targets: vec![],
            name: None,
            selector: None,
        };

        export_environments(output_dir.clone(), paths, false, opts)?;

        check_files(
            &output_dir,
            &[
                &format!("{}/static/initial-deployment.yaml", output_dir.display()),
                &format!("{}/static/initial-service.yaml", output_dir.display()),
                &format!("{}/manifest.json", output_dir.display()),
            ],
        )?;

        let manifest_content = fs::read_to_string(&manifest_path)?;
        let manifest: HashMap<String, String> = serde_json::from_str(&manifest_content)?;

        let mut expected_manifest = HashMap::new();
        expected_manifest.insert(
            "static/initial-deployment.yaml".to_string(),
            "test-export-envs/static-env/main.jsonnet".to_string(),
        );
        expected_manifest.insert(
            "static/initial-service.yaml".to_string(),
            "test-export-envs/static-env/main.jsonnet".to_string(),
        );

        assert_eq!(
            manifest, expected_manifest,
            "Manifest content doesn't match after deletion"
        );

        Ok(())
    })();

    // Always restore the original directory
    std::env::set_current_dir(original_dir)?;

    result
}

#[test]
fn test_export_environments_skip_manifest() -> Result<()> {
    let temp_dir = tempfile::tempdir()?;
    let output_dir = temp_dir.path().to_path_buf();

    let original_dir = std::env::current_dir()?;
    let testdata_dir = original_dir.join("testdata").canonicalize()?;
    std::env::set_current_dir(&testdata_dir)?;

    let result = (|| -> Result<()> {
        let paths = vec![PathBuf::from("test-export-envs")];
        let opts = ExportOptions {
            format: "{{.metadata.namespace}}/{{.metadata.name}}".to_string(),
            extension: "yaml".to_string(),
            parallelism: 1,
            cache_path: None,
            cache_path_regexes: vec![],
            merge_strategy: ExportMergeStrategy::None,
            merge_deleted_envs: vec![],
            skip_manifest: true,
            targets: vec![],
            name: None,
            selector: None,
        };

        export_environments(output_dir.clone(), paths, true, opts)?;

        // Check that manifest files exist but manifest.json does NOT exist
        check_files(
            &output_dir,
            &[
                &format!(
                    "{}/inline-namespace1/my-configmap.yaml",
                    output_dir.display()
                ),
                &format!(
                    "{}/inline-namespace1/my-deployment.yaml",
                    output_dir.display()
                ),
                &format!("{}/inline-namespace1/my-service.yaml", output_dir.display()),
                &format!(
                    "{}/inline-namespace2/my-deployment.yaml",
                    output_dir.display()
                ),
                &format!("{}/inline-namespace2/my-service.yaml", output_dir.display()),
                &format!("{}/static/initial-deployment.yaml", output_dir.display()),
                &format!("{}/static/initial-service.yaml", output_dir.display()),
            ],
        )?;

        // Verify manifest.json does not exist
        let manifest_path = output_dir.join("manifest.json");
        assert!(
            !manifest_path.exists(),
            "manifest.json should not exist when skip_manifest is true"
        );

        Ok(())
    })();

    std::env::set_current_dir(original_dir)?;

    result
}
