package jsonnet

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// To be consistent with the jsonnet executable,
// when evaluating a file, `std.thisFile` should point to the given path
func TestEvaluateFile(t *testing.T) {
	result, err := EvaluateFile("testdata/thisFile/main.jsonnet", Opts{})
	assert.NoError(t, err)
	assert.Equal(t, `{
   "test": "testdata/thisFile/main.jsonnet"
}
`, result)
}

func TestEvaluateFileDoesntExist(t *testing.T) {
	result, err := EvaluateFile("testdata/doesnt-exist/main.jsonnet", Opts{})
	assert.EqualError(t, err, "open testdata/doesnt-exist/main.jsonnet: no such file or directory")
	assert.Equal(t, "", result)
}
