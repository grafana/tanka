package tanka

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/labels"
)

func Test_replaceTmplText(t *testing.T) {
	type args struct {
		s   string
		old string
		new string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"text only", args{"a", "a", "b"}, "b"},
		{"action blocks", args{"{{a}}{{.}}", "a", "b"}, "{{a}}{{.}}"},
		{"mixed", args{"a{{a}}a{{a}}a", "a", "b"}, "b{{a}}b{{a}}b"},
		{"invalid template format handled as text", args{"a}}a{{a", "a", "b"}, "b}}b{{b"},
		{
			name: "keep path separator in action block",
			args: args{`{{index .metadata.labels "app.kubernetes.io/name"}}/{{.metadata.name}}`, "/", BelRune},
			want: "{{index .metadata.labels \"app.kubernetes.io/name\"}}\u0007{{.metadata.name}}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := replaceTmplText(tt.args.s, tt.args.old, tt.args.new); got != tt.want {
				t.Errorf("replaceInTmplText() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExportEnvironments(t *testing.T) {
	tempDir := t.TempDir()
	require.NoError(t, os.Chdir("testdata"))
	defer func() { require.NoError(t, os.Chdir("..")) }()

	// Find envs
	envs, err := FindEnvs(t.Context(), "test-export-envs", FindOpts{Selector: labels.Everything()})
	require.NoError(t, err)

	// Export all envs
	opts := &ExportEnvOpts{
		Format:    "{{.metadata.namespace}}/{{.metadata.name}}",
		Extension: "yaml",
	}
	opts.Opts.ExtCode = jsonnet.InjectedCode{
		"deploymentName": "'initial-deployment'",
		"serviceName":    "'initial-service'",
	}
	require.NoError(t, ExportEnvironments(t.Context(), envs, tempDir, opts))
	checkFiles(t, tempDir, []string{
		filepath.Join(tempDir, "inline-namespace1", "my-configmap.yaml"),
		filepath.Join(tempDir, "inline-namespace1", "my-deployment.yaml"),
		filepath.Join(tempDir, "inline-namespace1", "my-service.yaml"),
		filepath.Join(tempDir, "inline-namespace2", "my-deployment.yaml"),
		filepath.Join(tempDir, "inline-namespace2", "my-service.yaml"),
		filepath.Join(tempDir, "static", "initial-deployment.yaml"),
		filepath.Join(tempDir, "static", "initial-service.yaml"),
		filepath.Join(tempDir, "manifest.json"),
	})
	manifestContent, err := os.ReadFile(filepath.Join(tempDir, "manifest.json"))
	require.NoError(t, err)
	assert.Equal(t, string(manifestContent), `{
    "inline-namespace1/my-configmap.yaml": "test-export-envs/inline-envs/main.jsonnet",
    "inline-namespace1/my-deployment.yaml": "test-export-envs/inline-envs/main.jsonnet",
    "inline-namespace1/my-service.yaml": "test-export-envs/inline-envs/main.jsonnet",
    "inline-namespace2/my-deployment.yaml": "test-export-envs/inline-envs/main.jsonnet",
    "inline-namespace2/my-service.yaml": "test-export-envs/inline-envs/main.jsonnet",
    "static/initial-deployment.yaml": "test-export-envs/static-env/main.jsonnet",
    "static/initial-service.yaml": "test-export-envs/static-env/main.jsonnet"
}`)

	// Try to re-export
	assert.EqualError(t, ExportEnvironments(t.Context(), envs, tempDir, opts), fmt.Sprintf("output dir `%s` not empty. Pass a different --merge-strategy to ignore this", tempDir))

	// Try to re-export with the --merge-strategy=fail-on-conflicts flag. Will still fail because Tanka will not overwrite manifests silently
	opts.MergeStrategy = ExportMergeStrategyFailConflicts
	assert.ErrorContains(t, ExportEnvironments(t.Context(), envs, tempDir, opts), "already exists. Aborting")

	// Re-export only one env with --merge-stategy=replace-envs flag
	opts.Opts.ExtCode = jsonnet.InjectedCode{
		"deploymentName": "'updated-deployment'",
		"serviceName":    "'updated-service'",
	}
	opts.MergeStrategy = ExportMergeStrategyReplaceEnvs
	staticEnv, err := FindEnvs(t.Context(), "test-export-envs", FindOpts{Selector: labels.SelectorFromSet(labels.Set{"type": "static"})})
	require.NoError(t, err)
	require.NoError(t, ExportEnvironments(t.Context(), staticEnv, tempDir, opts))
	checkFiles(t, tempDir, []string{
		filepath.Join(tempDir, "inline-namespace1", "my-configmap.yaml"),
		filepath.Join(tempDir, "inline-namespace1", "my-deployment.yaml"),
		filepath.Join(tempDir, "inline-namespace1", "my-service.yaml"),
		filepath.Join(tempDir, "inline-namespace2", "my-deployment.yaml"),
		filepath.Join(tempDir, "inline-namespace2", "my-service.yaml"),
		filepath.Join(tempDir, "static", "updated-deployment.yaml"),
		filepath.Join(tempDir, "static", "updated-service.yaml"),
		filepath.Join(tempDir, "manifest.json"),
	})
	manifestContent, err = os.ReadFile(filepath.Join(tempDir, "manifest.json"))
	require.NoError(t, err)
	assert.Equal(t, string(manifestContent), `{
    "inline-namespace1/my-configmap.yaml": "test-export-envs/inline-envs/main.jsonnet",
    "inline-namespace1/my-deployment.yaml": "test-export-envs/inline-envs/main.jsonnet",
    "inline-namespace1/my-service.yaml": "test-export-envs/inline-envs/main.jsonnet",
    "inline-namespace2/my-deployment.yaml": "test-export-envs/inline-envs/main.jsonnet",
    "inline-namespace2/my-service.yaml": "test-export-envs/inline-envs/main.jsonnet",
    "static/updated-deployment.yaml": "test-export-envs/static-env/main.jsonnet",
    "static/updated-service.yaml": "test-export-envs/static-env/main.jsonnet"
}`)

	// Re-export and delete the files of one env
	opts.Opts.ExtCode = jsonnet.InjectedCode{
		"deploymentName": "'updated-again-deployment'",
		"serviceName":    "'updated-again-service'",
	}
	opts.MergeDeletedEnvs = []string{"test-export-envs/inline-envs/main.jsonnet"}
	require.NoError(t, ExportEnvironments(t.Context(), staticEnv, tempDir, opts))
	checkFiles(t, tempDir, []string{
		filepath.Join(tempDir, "static", "updated-again-deployment.yaml"),
		filepath.Join(tempDir, "static", "updated-again-service.yaml"),
		filepath.Join(tempDir, "manifest.json"),
	})
	manifestContent, err = os.ReadFile(filepath.Join(tempDir, "manifest.json"))
	require.NoError(t, err)
	assert.Equal(t, string(manifestContent), `{
    "static/updated-again-deployment.yaml": "test-export-envs/static-env/main.jsonnet",
    "static/updated-again-service.yaml": "test-export-envs/static-env/main.jsonnet"
}`)
}

func TestExportEnvironmentsBroken(t *testing.T) {
	tempDir := t.TempDir()
	require.NoError(t, os.Chdir("testdata"))
	defer func() { require.NoError(t, os.Chdir("..")) }()

	// Find envs
	envs, err := FindEnvs(t.Context(), "test-export-envs-broken", FindOpts{Selector: labels.Everything()})
	require.NoError(t, err)

	// Export all envs
	opts := &ExportEnvOpts{
		Format:    "{{.metadata.namespace}}/{{.metadata.name}}",
		Extension: "yaml",
	}

	var schemaError *manifest.SchemaError
	require.ErrorAs(t, ExportEnvironments(t.Context(), envs, tempDir, opts), &schemaError)
}

func BenchmarkExportEnvironmentsWithReplaceEnvs(b *testing.B) {
	zerolog.SetGlobalLevel(zerolog.Disabled)
	tempDir := b.TempDir()
	require.NoError(b, os.Chdir("testdata"))
	defer func() { require.NoError(b, os.Chdir("..")) }()

	// Find envs
	envs, err := FindEnvs(b.Context(), "test-export-envs", FindOpts{Selector: labels.Everything()})
	require.NoError(b, err)

	// Export all envs
	opts := &ExportEnvOpts{
		Format:        "{{.metadata.namespace}}/{{.metadata.name}}",
		Extension:     "yaml",
		MergeStrategy: ExportMergeStrategyReplaceEnvs,
	}
	opts.Opts.ExtCode = jsonnet.InjectedCode{
		"deploymentName": "'initial-deployment'",
		"serviceName":    "'initial-service'",
	}
	// Export a first time so that the benchmark loops are identical
	require.NoError(b, ExportEnvironments(b.Context(), envs, tempDir, opts))

	// On every loop, delete manifests from previous envs + reexport all envs
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		require.NoError(b, ExportEnvironments(b.Context(), envs, tempDir, opts), "failed on iteration %d", i)
	}
}

func TestExportEnvironmentsSkipManifest(t *testing.T) {
	tempDir := t.TempDir()
	require.NoError(t, os.Chdir("testdata"))
	defer func() { require.NoError(t, os.Chdir("..")) }()

	// Find envs
	envs, err := FindEnvs(t.Context(), "test-export-envs", FindOpts{Selector: labels.Everything()})
	require.NoError(t, err)

	// Export all envs with skip manifest flag
	opts := &ExportEnvOpts{
		Format:       "{{.metadata.namespace}}/{{.metadata.name}}",
		Extension:    "yaml",
		SkipManifest: true,
	}
	opts.Opts.ExtCode = jsonnet.InjectedCode{
		"deploymentName": "'test-deployment'",
		"serviceName":    "'test-service'",
	}
	require.NoError(t, ExportEnvironments(t.Context(), envs, tempDir, opts))
	
	// Check that all manifest files are created but manifest.json is NOT created
	expectedFiles := []string{
		filepath.Join(tempDir, "inline-namespace1", "my-configmap.yaml"),
		filepath.Join(tempDir, "inline-namespace1", "my-deployment.yaml"),
		filepath.Join(tempDir, "inline-namespace1", "my-service.yaml"),
		filepath.Join(tempDir, "inline-namespace2", "my-deployment.yaml"),
		filepath.Join(tempDir, "inline-namespace2", "my-service.yaml"),
		filepath.Join(tempDir, "static", "test-deployment.yaml"),
		filepath.Join(tempDir, "static", "test-service.yaml"),
	}
	checkFiles(t, tempDir, expectedFiles)
	
	// Explicitly verify manifest.json does not exist
	manifestPath := filepath.Join(tempDir, "manifest.json")
	_, err = os.Stat(manifestPath)
	assert.True(t, os.IsNotExist(err), "manifest.json should not exist when SkipManifest is true")
}

func checkFiles(t testing.TB, dir string, files []string) {
	t.Helper()

	var existingFiles []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		existingFiles = append(existingFiles, path)
		return nil
	})
	require.NoError(t, err)

	assert.ElementsMatch(t, files, existingFiles)
}
