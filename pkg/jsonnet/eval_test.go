package jsonnet

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/grafana/tanka/pkg/jsonnet/implementations/binary"
	"github.com/grafana/tanka/pkg/jsonnet/implementations/goimpl"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var jsonnetImpl = &goimpl.JsonnetGoImplementation{}

const importTreeResult = `[
   {
      "breed": "apple",
      "color": "red",
      "creates": "o2",
      "eats": "co2",
      "keeps": "the world healthy",
      "kind": "tree",
      "needs": "water",
      "size": "m"
   },
   {
      "breed": "cherry",
      "color": "red",
      "creates": "o2",
      "eats": "co2",
      "keeps": "the world healthy",
      "kind": "tree",
      "needs": "water",
      "size": "xs"
   },
   {
      "breed": "peach",
      "color": "orange",
      "creates": "o2",
      "eats": "co2",
      "keeps": "the world healthy",
      "kind": "tree",
      "needs": "water",
      "size": "s"
   }
]
`

const thisFileResult = `{
   "test": "testdata/thisFile/main.jsonnet"
}
`

// To be consistent with the jsonnet executable,
// when evaluating a file, `std.thisFile` should point to the given path
func TestEvaluateFile(t *testing.T) {
	result, err := EvaluateFile(jsonnetImpl, "testdata/thisFile/main.jsonnet", Opts{})
	assert.NoError(t, err)
	assert.Equal(t, thisFileResult, result)
}

func TestEvaluateFileWithInvalidBinary(t *testing.T) {
	binaryImpl := &binary.JsonnetBinaryImplementation{BinPath: "this-file-doesnt-exist"}
	result, err := EvaluateFile(binaryImpl, "testdata/thisFile/main.jsonnet", Opts{})
	assert.Equal(t, result, "")
	assert.ErrorIs(t, err, exec.ErrNotFound)
}

// This test requires jsonnet to be installed and available in the PATH
func TestEvaluateFileWithJsonnetBinary(t *testing.T) {
	binaryImpl := &binary.JsonnetBinaryImplementation{BinPath: "jsonnet"}
	result, err := EvaluateFile(binaryImpl, "testdata/thisFile/main.jsonnet", Opts{})
	assert.NoError(t, err)
	assert.Equal(t, thisFileResult, result)
}

func TestEvaluateFileDoesntExist(t *testing.T) {
	result, err := EvaluateFile(jsonnetImpl, "testdata/doesnt-exist/main.jsonnet", Opts{})
	assert.EqualError(t, err, "open testdata/doesnt-exist/main.jsonnet: no such file or directory")
	assert.Equal(t, "", result)
}

func TestEvaluateFileWithCaching(t *testing.T) {
	tmp, err := os.MkdirTemp("", "test-tanka-caching")
	require.NoError(t, err)
	defer os.RemoveAll(tmp)
	cachePath := filepath.Join(tmp, "cache") // Should be created during caching

	// Evaluate two files
	result, err := EvaluateFile(jsonnetImpl, "testdata/thisFile/main.jsonnet", Opts{CachePath: cachePath})
	assert.NoError(t, err)
	assert.Equal(t, thisFileResult, result)
	result, err = EvaluateFile(jsonnetImpl, "testdata/importTree/main.jsonnet", Opts{CachePath: cachePath})
	assert.NoError(t, err)
	assert.Equal(t, importTreeResult, result)

	// Check that we have two entries in the cache
	readCache, err := os.ReadDir(cachePath)
	require.NoError(t, err)
	assert.Len(t, readCache, 2)

	// Evaluate two files again, same result
	result, err = EvaluateFile(jsonnetImpl, "testdata/thisFile/main.jsonnet", Opts{CachePath: cachePath})
	assert.NoError(t, err)
	assert.Equal(t, thisFileResult, result)
	result, err = EvaluateFile(jsonnetImpl, "testdata/importTree/main.jsonnet", Opts{CachePath: cachePath})
	assert.NoError(t, err)
	assert.Equal(t, importTreeResult, result)

	// Modify the cache items
	for _, entry := range readCache {
		require.NoError(t, os.WriteFile(filepath.Join(cachePath, entry.Name()), []byte(entry.Name()), 0666))
	}

	// Evaluate two files again, modified cache is returned instead of the actual result
	result, err = EvaluateFile(jsonnetImpl, "testdata/thisFile/main.jsonnet", Opts{CachePath: cachePath})
	assert.NoError(t, err)
	assert.Equal(t, "BYfdlr1ZOVwiOfbd89JYTcK-eRQh05bi8ky3k1vVW5o=.json", result)
	result, err = EvaluateFile(jsonnetImpl, "testdata/importTree/main.jsonnet", Opts{CachePath: cachePath})
	assert.NoError(t, err)
	assert.Equal(t, "R_3hy-dRfOwXN-fezQ50ZF4dnrFcBcbQ9LztR_XWzJA=.json", result)
}
