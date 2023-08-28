package tanka

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkFindEnvsFromSinglePath(b *testing.B) {
	tempDir, envPaths := buildLargeEnvironmentDirForFindTest(b)
	require.Len(b, envPaths, 105) // 100 static + 5 inline

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		envs, err := FindEnvs(context.Background(), tempDir, FindOpts{})
		require.Len(b, envs, 200)
		require.NoError(b, err)
	}
}

func BenchmarkFindEnvsFromPaths(b *testing.B) {
	_, envPaths := buildLargeEnvironmentDirForFindTest(b)
	require.Len(b, envPaths, 105) // 100 static + 5 inline

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		envs, err := FindEnvsFromPaths(context.Background(), envPaths, FindOpts{})
		require.Len(b, envs, 200)
		require.NoError(b, err)
	}
}

// create a directory with lots of inline and static environments
func buildLargeEnvironmentDirForFindTest(t testing.TB) (string, []string) {
	t.Helper()

	// create a temp dir
	tempDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(tempDir, "jsonnetfile.json"), []byte(`{}`), 0644))
	var envPaths []string

	// create 100 static envs: 5 indent levels (dir-0, dir-0/dir-1, etc) + 20 static envs each)
	// create 100 inline envs: 5 indent levels (dir-0, dir-0/dir-1, etc) + 1 inline env dir in each with 20 envs each)
	for indent := range [5]struct{}{} {
		// create indented dir
		envDir := tempDir
		for i := 0; i < indent; i++ {
			envDir = filepath.Join(envDir, fmt.Sprintf("dir-%d", indent))
		}

		// create 20 static envs
		for id := range [20]struct{}{} {
			staticEnvDir := filepath.Join(envDir, fmt.Sprintf("static-%d", id))
			staticEnvFile := filepath.Join(staticEnvDir, "main.jsonnet")
			envPaths = append(envPaths, staticEnvDir)
			require.NoError(t, os.MkdirAll(staticEnvDir, 0755))
			require.NoError(t, os.WriteFile(staticEnvFile, []byte(`{}`), 0644))
			require.NoError(t, os.WriteFile(filepath.Join(staticEnvDir, "spec.json"), []byte(fmt.Sprintf(`{
    "apiVersion": "tanka.dev/v1alpha1",
    "kind": "Environment",
    "metadata": {
		"name": "%[1]s",
		"labels": {}
    },
    "spec": {
		"apiServer": "https://192.168.0.1",
		"namespace": "blabla",
		"cluster": "blabla"
	}
}`, staticEnvDir)), 0644))
		}

		inlineEnvDir := filepath.Join(envDir, "inline")
		inlineEnvFile := filepath.Join(inlineEnvDir, "main.jsonnet")
		envPaths = append(envPaths, inlineEnvDir)
		require.NoError(t, os.MkdirAll(inlineEnvDir, 0755))
		require.NoError(t, os.WriteFile(inlineEnvFile, []byte(fmt.Sprintf(`[
{
	"apiVersion": "tanka.dev/v1alpha1",
	"kind": "Environment",
	"metadata": {
		"name": "%[1]s/%%s",
		"labels": {}
	},
	"spec": {
		"apiServer": "https://192.168.0.1",
		"namespace": "blabla",
		"cluster": "blabla",
	}
} for i in std.range(0, 19)
		]`, inlineEnvDir)), 0644))
	}

	return tempDir, envPaths
}
