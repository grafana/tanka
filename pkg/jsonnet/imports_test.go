package jsonnet

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/grafana/tanka/pkg/jsonnet/implementations/goimpl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTransitiveImports checks that TransitiveImports is able to report all
// recursive imports of a file
func TestTransitiveImports(t *testing.T) {
	imports, err := TransitiveImports("testdata/importTree")
	fmt.Println(imports)
	require.NoError(t, err)
	assert.Equal(t, []string{
		"main.jsonnet",
		"trees.jsonnet",
		"trees/apple.jsonnet",
		"trees/cherry.jsonnet",
		"trees/generic.libsonnet",
		"trees/peach.jsonnet",
	}, imports)
}

func BenchmarkGetSnippetHash(b *testing.B) {
	for _, tc := range []struct {
		name           string
		importFromMain bool
		expectedHash   string
	}{
		{
			name:           "all-imported-from-main",
			importFromMain: true,
			expectedHash:   "ktY8NYZOoPacsNYrH7-DslRgLG54EMRdk3MQSM3vcUg=",
		},
		{
			name:           "deeply-nested",
			importFromMain: false,
			expectedHash:   "W1Q_uS6jTGcsd7nvJfc-i785sqjBmOzfOAzqzhXVc0A=",
		},
	} {
		b.Run(tc.name, func(b *testing.B) {
			// Create a very large and complex project
			tempDir := b.TempDir()
			generateTestProject(b, tempDir, 1000, tc.importFromMain)

			// Create a VM. It's important to reuse the same VM
			// While there is a caching mechanism that normally shouldn't be shared in a benchmark iteration,
			// it's useful to evaluate its impact here, because the caching will also improve the evaluation performance afterwards.
			vm := goimpl.MakeRawVM([]string{tempDir}, nil, nil, 0, false, false)
			content, err := os.ReadFile(filepath.Join(tempDir, "main.jsonnet"))
			require.NoError(b, err)

			// Run the benchmark
			mainPath := filepath.Join(tempDir, "main.jsonnet")
			c := string(content)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				fileHashes = sync.Map{}
				hash, err := getSnippetHash(vm, mainPath, c)
				require.NoError(b, err)
				require.Equal(b, tc.expectedHash, hash)
			}
		})
	}
}

func generateTestProject(t testing.TB, dir string, depth int, importAllFromMain bool) []string {
	t.Helper()

	const testFile = `
	local localImport = <IMPORT>;
	local myFunc = function() <IMPORT>;
	
	{
		local this = self,
	
		attribute: {
			name: 'test',
			value: self.name,
			otherValue: 'other ' + self.value,
		},
		nested: {
			nested: {
				nested: {
					nested: {
						nested1: {
							nested: {
								nested1: {
									nested: {
										attribute: <IMPORT>,
									},
								},
								nested2: {
									strValue: this.nested.nested.nested,
								},
							},
						},
						nested2: {
							intValue: 1,
							importValue: <IMPORT>,
						},
					},
				},
			},
		},
	
		other: myFunc(),
		useLocal: localImport,
	}`

	var allFiles []string

	var mainContentSplit []string
	for i := 0; i < depth; i++ {
		mainContentSplit = append(mainContentSplit, fmt.Sprintf("(import 'file%d.libsonnet')", i))
		filePath := filepath.Join(dir, fmt.Sprintf("file%d.libsonnet", i))
		err := os.WriteFile(
			filePath,
			[]byte(strings.ReplaceAll(testFile, "<IMPORT>", fmt.Sprintf("import 'file%d.libsonnet'", i+1))),
			0644,
		)
		require.NoError(t, err)
		allFiles = append(allFiles, filePath)
	}
	if !importAllFromMain {
		mainContentSplit = append(mainContentSplit, "import 'file0.libsonnet'")
	}
	require.NoError(t, os.WriteFile(filepath.Join(dir, "main.jsonnet"), []byte(strings.Join(mainContentSplit, " + ")), 0644))
	allFiles = append(allFiles, filepath.Join(dir, "main.jsonnet"))
	require.NoError(t, os.WriteFile(filepath.Join(dir, fmt.Sprintf("file%d.libsonnet", depth)), []byte(`"a string"`), 0644))
	allFiles = append(allFiles, filepath.Join(dir, fmt.Sprintf("file%d.libsonnet", depth)))

	return allFiles
}
