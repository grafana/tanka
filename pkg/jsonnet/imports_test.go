package jsonnet

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func BenchmarkTransitiveImports(b *testing.B) {
	imports, err := TransitiveImports("testdata/importTree")
	require.NoError(b, err)
	assert.Equal(b, []string{
		"main.jsonnet",
		"trees.jsonnet",
		"trees/apple.jsonnet",
		"trees/cherry.jsonnet",
		"trees/generic.libsonnet",
		"trees/leaf.libsonnet",
		"trees/peach.jsonnet",
	}, imports)
}

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
		"trees/leaf.libsonnet",
		"trees/peach.jsonnet",
	}, imports)
}
