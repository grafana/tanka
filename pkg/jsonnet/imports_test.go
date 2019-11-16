package jsonnet

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTransitiveImports checks that TransitiveImports is able to report all
// recursive imports of a file
func TestTransitiveImports(t *testing.T) {
	imports, err := TransitiveImports("testdata/main.jsonnet")
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{
		"testdata/trees.jsonnet",

		"testdata/trees/apple.jsonnet",
		"testdata/trees/cherry.jsonnet",
		"testdata/trees/peach.jsonnet",

		"testdata/trees/generic.libsonnet",
	}, imports)
}
