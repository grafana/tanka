package jsonnet

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/grafana/tanka/pkg/jsonnet/jpath"
	"github.com/stretchr/testify/require"
)

type findImportersTestCase struct {
	name              string
	files             []string
	expectedImporters []string
	expectedErr       error

	checkTransitiveImporters    bool
	expectedTransitiveImporters []string
}

func (tc findImportersTestCase) run(t testing.TB) {
	importers, err := FindImporterForFiles(t.Context(), "testdata/findImporters", tc.files)

	if tc.expectedErr != nil {
		require.EqualError(t, err, tc.expectedErr.Error())
	} else {
		require.NoError(t, err)
		require.Equal(t, tc.expectedImporters, importers)

		if tc.checkTransitiveImporters {
			transitiveImporters, err := FindTransitiveImportersForFile("testdata/findImporters", tc.files)
			require.NoError(t, err)
			require.Equal(t, tc.expectedTransitiveImporters, transitiveImporters)
		}
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
			expectedTransitiveImporters: []string{
				absPath(t, "testdata/findImporters/environments/imports-lib-and-vendored-through-chain/chain1.libsonnet"),
				absPath(t, "testdata/findImporters/environments/imports-lib-and-vendored-through-chain/chain2.libsonnet"),
				absPath(t, "testdata/findImporters/environments/imports-lib-and-vendored-through-chain/main.jsonnet"),
				absPath(t, "testdata/findImporters/environments/imports-locals-and-vendored/main.jsonnet"),
				absPath(t, "testdata/findImporters/environments/imports-symlinked-vendor/main.jsonnet"),
				absPath(t, "testdata/findImporters/lib/lib1/main.libsonnet"),
				absPath(t, "testdata/findImporters/vendor/vendor-symlinked/main.libsonnet"),
				absPath(t, "testdata/findImporters/vendor/vendored/main.libsonnet"),
			},
			checkTransitiveImporters: true,
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
		{
			name:  "imported file in lib relative to env",
			files: []string{"testdata/findImporters/environments/lib-import-relative-to-env/file-to-import.libsonnet"},
			expectedImporters: []string{
				absPath(t, "testdata/findImporters/environments/lib-import-relative-to-env/folder1/folder2/main.jsonnet"),
			},
		},
		{
			name: "unused deleted file",
			files: []string{
				"deleted:testdata/findImporters/vendor/deleted-vendor/main.libsonnet",
			},
			expectedImporters: nil,
		},
		{
			name: "deleted local path that is still potentially imported",
			files: []string{
				"deleted:testdata/findImporters/environments/using-deleted-stuff/my-import-dir/main.libsonnet",
			},
			expectedImporters: []string{
				absPath(t, "testdata/findImporters/environments/using-deleted-stuff/main.jsonnet"),
			},
		},
		{
			name: "deleted lib that is still potentially imported",
			files: []string{
				"deleted:testdata/findImporters/lib/my-import-dir/main.libsonnet",
			},
			expectedImporters: []string{
				absPath(t, "testdata/findImporters/environments/using-deleted-stuff/main.jsonnet"),
			},
		},
		{
			name: "deleted vendor that is still potentially imported",
			files: []string{
				"deleted:testdata/findImporters/vendor/my-import-dir/main.libsonnet",
			},
			expectedImporters: []string{
				absPath(t, "testdata/findImporters/environments/using-deleted-stuff/main.jsonnet"),
			},
		},
		{
			name: "deleted lib that is still potentially imported, relative path from root",
			files: []string{
				"deleted:lib/my-import-dir/main.libsonnet",
			},
			expectedImporters: []string{
				absPath(t, "testdata/findImporters/environments/using-deleted-stuff/main.jsonnet"),
			},
		},
		{
			// All files in an environment are considered to be imported by the main file, so the same should apply for deleted files
			name: "deleted dir in environment",
			files: []string{
				"deleted:testdata/findImporters/environments/no-imports/deleted-dir/deleted-file.libsonnet",
			},
			expectedImporters: []string{
				absPath(t, "testdata/findImporters/environments/no-imports/main.jsonnet"),
			},
		},
		{
			name: "imports through a main file are followed",
			files: []string{
				"testdata/findImporters/environments/import-other-main-file/env2/file.libsonnet",
			},
			expectedImporters: []string{
				absPath(t, "testdata/findImporters/environments/import-other-main-file/env1/main.jsonnet"),
				absPath(t, "testdata/findImporters/environments/import-other-main-file/env2/main.jsonnet"),
			},
		},
	}
}

func TestFindImportersForFiles(t *testing.T) {
	// Sanity check
	// Make sure the main files all eval correctly
	// We want to make sure that the importers command works correctly,
	// but there's no point in testing on invalid jsonnet files
	files, err := FindFiles("testdata", nil)
	require.NoError(t, err)
	require.NotEmpty(t, files)
	for _, file := range files {
		// This project is known to be invalid (as the name suggests)
		if strings.Contains(file, "using-deleted-stuff") {
			continue
		}

		// Skip non-main files
		if filepath.Base(file) != jpath.DefaultEntrypoint {
			continue
		}
		_, err := EvaluateFile(t.Context(), jsonnetImpl, file, Opts{})
		require.NoError(t, err, "failed to eval %s", file)
	}

	for _, c := range findImportersTestCases(t) {
		t.Run(c.name, func(t *testing.T) {
			c.run(t)
		})
	}
}

func TestCountImporters(t *testing.T) {
	testcases := []struct {
		name       string
		dir        string
		recursive  bool
		fileRegexp string
		expected   string
	}{
		{
			name:      "project with no imports",
			dir:       "testdata/findImporters/environments/no-imports",
			recursive: true,
			expected:  "",
		},
		{
			name:      "project with imports",
			dir:       "testdata/findImporters/environments/imports-locals-and-vendored",
			recursive: true,
			expected: `testdata/findImporters/environments/imports-locals-and-vendored/local-file1.libsonnet: 1
testdata/findImporters/environments/imports-locals-and-vendored/local-file2.libsonnet: 1
`,
		},
		{
			name:      "lib non-recursive",
			dir:       "testdata/findImporters/lib/lib1",
			recursive: false,
			expected: `testdata/findImporters/lib/lib1/main.libsonnet: 1
`,
		},
		{
			name:      "lib recursive",
			dir:       "testdata/findImporters/lib/lib1",
			recursive: true,
			expected: `testdata/findImporters/lib/lib1/main.libsonnet: 1
testdata/findImporters/lib/lib1/subfolder/test.libsonnet: 0
`,
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			count, err := CountImporters(t.Context(), "testdata/findImporters", tc.dir, tc.recursive, tc.fileRegexp)
			require.NoError(t, err)
			require.Equal(t, tc.expected, count)
		})
	}
}

func BenchmarkFindImporters(b *testing.B) {
	// Create a very large and complex project
	tempDir, err := filepath.EvalSymlinks(b.TempDir())
	require.NoError(b, err)
	generateTestProject(b, tempDir, 100, false)

	// Run the benchmark
	expectedImporters := []string{filepath.Join(tempDir, "main.jsonnet")}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		importersCache = make(map[string][]string)
		jsonnetFilesCache = make(map[string]map[string]*cachedJsonnetFile)
		symlinkCache = make(map[string]string)
		importers, err := FindImporterForFiles(b.Context(), tempDir, []string{filepath.Join(tempDir, "file10.libsonnet")})

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
