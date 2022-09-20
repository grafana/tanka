package tanka

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/grafana/tanka/pkg/jsonnet"
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
	defer os.Chdir("..")

	// Find envs
	envs, err := FindEnvs("test-export-envs", FindOpts{Selector: labels.Everything()})
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
	require.NoError(t, ExportEnvironments(envs, tempDir, opts))
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
	assert.EqualError(t, ExportEnvironments(envs, tempDir, opts), fmt.Sprintf("Output dir `%s` not empty. Pass --merge to ignore this", tempDir))

	// Try to re-export with the --merge flag. Will still fail because Tanka will not overwrite manifests silently
	opts.Merge = true
	assert.ErrorContains(t, ExportEnvironments(envs, tempDir, opts), "already exists. Aborting")

	// Re-export only one env with --delete-previous flag
	opts.Opts.ExtCode = jsonnet.InjectedCode{
		"deploymentName": "'updated-deployment'",
		"serviceName":    "'updated-service'",
	}
	opts.DeletePrevious = true
	staticEnv, err := FindEnvs("test-export-envs", FindOpts{Selector: labels.SelectorFromSet(labels.Set{"type": "static"})})
	require.NoError(t, err)
	require.NoError(t, ExportEnvironments(staticEnv, tempDir, opts))
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
}

func BenchmarkExportEnvironmentsWithDeletePrevious(b *testing.B) {
	log.SetOutput(ioutil.Discard)
	tempDir := b.TempDir()
	require.NoError(b, os.Chdir("testdata"))
	defer os.Chdir("..")

	// Find envs
	envs, err := FindEnvs("test-export-envs", FindOpts{Selector: labels.Everything()})
	require.NoError(b, err)

	// Export all envs
	opts := &ExportEnvOpts{
		Format:         "{{.metadata.namespace}}/{{.metadata.name}}",
		Extension:      "yaml",
		Merge:          true,
		DeletePrevious: true,
	}
	opts.Opts.ExtCode = jsonnet.InjectedCode{
		"deploymentName": "'initial-deployment'",
		"serviceName":    "'initial-service'",
	}
	// Export a first time so that the benchmark loops are identical
	require.NoError(b, ExportEnvironments(envs, tempDir, opts))

	// On every loop, delete manifests from previous envs + reexport all envs
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		require.NoError(b, ExportEnvironments(envs, tempDir, opts), "failed on iteration %d", i)
	}
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
