package jsonnet

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

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

func BenchmarkGetSnippetHash(b *testing.B) {
	// Create a very large and complex project
	tempDir := b.TempDir()

	var mainContentSplit []string
	for i := 0; i < 1000; i++ {
		mainContentSplit = append(mainContentSplit, fmt.Sprintf("(import 'file%d.libsonnet')", i))
	}
	require.NoError(b, os.WriteFile(filepath.Join(tempDir, "main.jsonnet"), []byte(strings.Join(mainContentSplit, " + ")), 0644))
	for i := 0; i < 1000; i++ {
		err := os.WriteFile(
			filepath.Join(tempDir, fmt.Sprintf("file%d.libsonnet", i)),
			[]byte(strings.ReplaceAll(testFile, "<IMPORT>", fmt.Sprintf("import 'file%d.libsonnet'", i+1))),
			0644,
		)
		require.NoError(b, err)
	}
	require.NoError(b, os.WriteFile(filepath.Join(tempDir, "file1000.libsonnet"), []byte(`"a string"`), 0644))

	// Create a VM. It's important to reuse the same VM
	// While there is a caching mechanism that normally shouldn't be shared in a benchmark iteration,
	// it's useful to evaluate its impact here, because the caching will also improve the evaluation performance afterwards.
	vm := MakeVM(Opts{ImportPaths: []string{tempDir}})
	content, err := os.ReadFile(filepath.Join(tempDir, "main.jsonnet"))
	require.NoError(b, err)

	// Run the benchmark
	for i := 0; i < b.N; i++ {
		fileHashes = sync.Map{}
		hash, err := getSnippetHash(vm, filepath.Join(tempDir, "main.jsonnet"), string(content))
		require.NoError(b, err)
		require.Equal(b, "XrkW8N2EvkFMvdIuHTsGsQespVUl9_xiFmM7v1mqX5s=", hash)
	}
}
