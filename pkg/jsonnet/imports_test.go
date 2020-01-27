package jsonnet

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestTransitiveImports checks that TransitiveImports is able to report all
// recursive imports of a file
func TestTransitiveImports(t *testing.T) {
	imports, err := TransitiveImports("testdata")
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
