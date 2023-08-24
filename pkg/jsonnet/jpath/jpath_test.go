package jpath_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/jsonnet/implementations/goimpl"
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
