package jsonnet

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/stretchr/testify/require"
)

type findImportersTestCase struct {
	name              string
	files             []string
	expectedImporters []string
	expectedErr       error
}

func (tc findImportersTestCase) run(t testing.TB) {
	importers, err := FindImporterForFiles("testdata/findImporters", tc.files)

	if tc.expectedErr != nil {
		require.EqualError(t, err, tc.expectedErr.Error())
	} else {
		require.NoError(t, err)
		require.Equal(t, tc.expectedImporters, importers)
	}
}

func findImportersTestCases(t testing.TB) []findImportersTestCase {
	t.Helper()

	return []findImportersTestCase{
		{
			name:              "no files",
			files:             []string{},
			expectedImporters: nil,
		},
		{
			name:        "invalid file",
			files:       []string{"testdata/findImporters/does-not-exist.jsonnet"},
			expectedErr: fmt.Errorf("lstat %s: no such file or directory", absPath(t, "testdata/findImporters/does-not-exist.jsonnet")),
		},
		{
			name:              "project with no imports",
			files:             []string{"testdata/findImporters/environments/no-imports/main.jsonnet"},
			expectedImporters: []string{absPath(t, "testdata/findImporters/environments/no-imports/main.jsonnet")}, // itself only
		},
		{
			name:              "local import",
			files:             []string{"testdata/findImporters/environments/imports-locals-and-vendored/local-file1.libsonnet"},
			expectedImporters: []string{absPath(t, "testdata/findImporters/environments/imports-locals-and-vendored/main.jsonnet")},
		},
		{
			name:              "local import with relative path",
			files:             []string{"testdata/findImporters/environments/imports-locals-and-vendored/local-file2.libsonnet"},
			expectedImporters: []string{absPath(t, "testdata/findImporters/environments/imports-locals-and-vendored/main.jsonnet")},
		},
		{
			name:  "lib imported through chain",
			files: []string{"testdata/findImporters/lib/lib1/main.libsonnet"},
			expectedImporters: []string{
				absPath(t, "testdata/findImporters/environments/imports-lib-and-vendored-through-chain/main.jsonnet"),
			},
		},
		{
			name:  "vendored lib imported through chain + directly",
			files: []string{"testdata/findImporters/vendor/vendored/main.libsonnet"},
			expectedImporters: []string{
				absPath(t, "testdata/findImporters/environments/imports-lib-and-vendored-through-chain/main.jsonnet"),
				absPath(t, "testdata/findImporters/environments/imports-locals-and-vendored/main.jsonnet"),
				absPath(t, "testdata/findImporters/environments/imports-symlinked-vendor/main.jsonnet"),
			},
		},
		{
			name:  "vendored lib found through symlink", // expect same result as previous test
			files: []string{"testdata/findImporters/vendor/vendor-symlinked/main.libsonnet"},
			expectedImporters: []string{
				absPath(t, "testdata/findImporters/environments/imports-lib-and-vendored-through-chain/main.jsonnet"),
				absPath(t, "testdata/findImporters/environments/imports-locals-and-vendored/main.jsonnet"),
				absPath(t, "testdata/findImporters/environments/imports-symlinked-vendor/main.jsonnet"),
			},
		},
		{
			name:  "text file",
			files: []string{"testdata/findImporters/vendor/vendored/text-file.txt"},
			expectedImporters: []string{
				absPath(t, "testdata/findImporters/environments/imports-lib-and-vendored-through-chain/main.jsonnet"),
				absPath(t, "testdata/findImporters/environments/imports-locals-and-vendored/main.jsonnet"),
				absPath(t, "testdata/findImporters/environments/imports-symlinked-vendor/main.jsonnet"),
			},
		},
		{
			name:  "relative imported environment",
			files: []string{"testdata/findImporters/environments/relative-imported/main.jsonnet"},
			expectedImporters: []string{
				absPath(t, "testdata/findImporters/environments/relative-import/main.jsonnet"),
				absPath(t, "testdata/findImporters/environments/relative-imported/main.jsonnet"), // itself, it's a main file
			},
		},
		{
			name:  "relative imported environment with doubled '..'",
			files: []string{"testdata/findImporters/environments/relative-imported2/main.jsonnet"},
			expectedImporters: []string{
				absPath(t, "testdata/findImporters/environments/relative-import/main.jsonnet"),
				absPath(t, "testdata/findImporters/environments/relative-imported2/main.jsonnet"), // itself, it's a main file
			},
		},
		{
			name:  "relative imported text file",
			files: []string{"testdata/findImporters/other-files/test.txt"},
			expectedImporters: []string{
				absPath(t, "testdata/findImporters/environments/relative-import/main.jsonnet"),
			},
		},
		{
			name:  "relative imported text file with doubled '..'",
			files: []string{"testdata/findImporters/other-files/test2.txt"},
			expectedImporters: []string{
				absPath(t, "testdata/findImporters/environments/relative-import/main.jsonnet"),
			},
		},
		{
			name: "vendor override in env: override vendor used",
			files: []string{
				"testdata/findImporters/environments/vendor-override-in-env/vendor/vendor-override-in-env/main.libsonnet",
			},
			expectedImporters: []string{
				absPath(t, "testdata/findImporters/environments/vendor-override-in-env/main.jsonnet"),
			},
		},
		{
			name: "vendor override in env: global vendor unused",
			files: []string{
				"testdata/findImporters/vendor/vendor-override-in-env/main.libsonnet",
			},
			expectedImporters: nil,
		},
	}
}

func TestFindImportersForFiles(t *testing.T) {
	// Make sure the main files all eval correctly
	// We want to make sure that the importers command works correctly,
	// but there's no point in testing on invalid jsonnet files
	files, err := FindFiles("testdata", nil)
	require.NoError(t, err)
	require.NotEmpty(t, files)
	for _, file := range files {
		// Skip non-main files
		if filepath.Base(file) != jpath.DefaultEntrypoint {
			continue
		}
		_, err := EvaluateFile(file, Opts{})
		require.NoError(t, err, "failed to eval %s", file)
	}

	for _, c := range findImportersTestCases(t) {
		t.Run(c.name, func(t *testing.T) {
			c.run(t)
		})
	}
}

func BenchmarkFindImporters(b *testing.B) {
	// Create a very large and complex project
	tempDir := b.TempDir()
	generateTestProject(b, tempDir, 100, false)

	// Run the benchmark
	expectedImporters := []string{filepath.Join(tempDir, "main.jsonnet")}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		importersCache = make(map[string][]string)
		jsonnetFilesCache = make(map[string]map[string]*cachedJsonnetFile)
		symlinkCache = make(map[string]string)
		importers, err := FindImporterForFiles(tempDir, []string{filepath.Join(tempDir, "file10.libsonnet")})

		require.NoError(b, err)
		require.Equal(b, expectedImporters, importers)
	}
}

func BenchmarkFindImporters_StaticCases(b *testing.B) {
	for _, c := range findImportersTestCases(b) {
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				c.run(b)
			}
		})
	}
}

func absPath(t testing.TB, path string) string {
	t.Helper()

	abs, err := filepath.Abs(path)
	require.NoError(t, err)
	return abs
}
