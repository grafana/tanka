package jsonnet

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindImportersForFiles(t *testing.T) {
	cases := []struct {
		name              string
		files             []string
		expectedImporters []string
		expectedErr       error
	}{
		{
			name:  "no files",
			files: []string{},
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
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			importers, err := FindImporterForFiles("testdata/findImporters", c.files)

			if c.expectedErr != nil {
				require.EqualError(t, err, c.expectedErr.Error())
			} else {
				require.NoError(t, err)
				require.Equal(t, c.expectedImporters, importers)
			}
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
		importers, err := FindImporterForFiles(tempDir, []string{filepath.Join(tempDir, "file10.libsonnet")})

		require.NoError(b, err)
		require.Equal(b, expectedImporters, importers)
	}
}

func absPath(t *testing.T, path string) string {
	t.Helper()

	abs, err := filepath.Abs(path)
	require.NoError(t, err)
	return abs
}
