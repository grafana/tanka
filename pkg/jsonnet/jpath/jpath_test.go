package jpath_test

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/jsonnet/implementations/goimpl"
	"github.com/grafana/tanka/pkg/jsonnet/jpath"
)

var jsonnetImpl = &goimpl.JsonnetGoImplementation{}

func TestResolvePrecedence(t *testing.T) {
	s, err := jsonnet.EvaluateFile(jsonnetImpl, "./testdata/precedence/environments/default/main.jsonnet", jsonnet.Opts{})
	require.NoError(t, err)

	want := map[string]string{
		"baseDir":        "baseDir",
		"lib":            "/lib",
		"baseDir-vendor": "baseDir-vendor",
		"vendor":         "/vendor",
	}

	w, err := json.Marshal(want)
	require.NoError(t, err)

	assert.JSONEq(t, string(w), s)
}

func TestJPathWeights(t *testing.T) {
	path := "./testdata/valid/environments/default/main.jsonnet"

	t.Run("default-weights", func(t *testing.T) {
		paths, _, root, err := jpath.Resolve(path, false, nil)
		require.Equal(t, []string{
			filepath.Join(root, "vendor"),
			filepath.Join(root, "environments", "default", "vendor"),
			filepath.Join(root, "lib"),
			filepath.Join(root, "environments", "default"),
		}, paths)
		require.NoError(t, err)
	})

	t.Run("custom-paths", func(t *testing.T) {
		_, _, root, err := jpath.Resolve(path, false, nil)
		require.NoError(t, err)
		paths, _, root, err := jpath.Resolve(path, false, []jpath.WeightedJPath{
			jpath.NewStaticallyWeightedJPath(filepath.Join(root, "vendor-dev"), 250),
			jpath.NewStaticallyWeightedJPath(filepath.Join(root, "prio-lib"), 2),
		})
		require.Equal(t, []string{
			filepath.Join(root, "vendor"),
			filepath.Join(root, "vendor-dev"),
			filepath.Join(root, "environments", "default", "vendor"),
			filepath.Join(root, "lib"),
			filepath.Join(root, "prio-lib"),
			filepath.Join(root, "environments", "default"),
		}, paths)
		require.NoError(t, err)
	})
}
