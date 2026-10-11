package native

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseYAMLDocuments(t *testing.T) {
	const input = `# Comments do not create a document.
apiVersion: v1
kind: ConfigMap
metadata:
  name: example
data:
  boolean: "true"
  number: "001"
  "null": "null"
  script: |
    echo first
    echo second
---
# Preserve document order and JSON-compatible types in nested values.
defaults: &defaults
  enabled: true
  replicas: 3
config:
  <<: *defaults
  replicas: 5
  ports: [80, 443]
  optional: null
`
	want := []interface{}{
		map[string]interface{}{
			"apiVersion": "v1", "kind": "ConfigMap",
			"metadata": map[string]interface{}{"name": "example"},
			"data": map[string]interface{}{
				"boolean": "true", "number": "001", "null": "null",
				"script": "echo first\necho second\n",
			},
		},
		map[string]interface{}{
			"defaults": map[string]interface{}{"enabled": true, "replicas": float64(3)},
			"config": map[string]interface{}{
				"enabled": true, "replicas": float64(5),
				"ports": []interface{}{float64(80), float64(443)}, "optional": nil,
			},
		},
	}

	got, err, callErr := callNative("parseYaml", []interface{}{input})
	require.NoError(t, callErr)
	require.NoError(t, err)
	require.Equal(t, want, got)
}

func TestParseYAMLDocumentErrors(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
	}{
		{name: "invalid later document", input: "valid: true\n---\ninvalid: ["},
		{name: "unknown alias", input: "value: *missing"},
		{name: "non-JSON number", input: "value: .nan"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err, callErr := callNative("parseYaml", []interface{}{tc.input})
			require.NoError(t, callErr)
			require.Error(t, err)
			require.Nil(t, got, "do not return a partial document stream")
		})
	}
}

func TestManifestYAMLFromJSONScalarTypes(t *testing.T) {
	const input = `{
  "data": {
    "boolean": "true",
    "null": "null",
    "number": "001",
    "script": "echo first\necho second\n"
  },
  "spec": {
    "enabled": true,
    "optional": null,
    "ports": [80, 443],
    "replicas": 3
  }
}`
	const want = `data:
    boolean: "true"
    "null": "null"
    number: "001"
    script: |
        echo first
        echo second
spec:
    enabled: true
    optional: null
    ports:
        - 80
        - 443
    replicas: 3
`

	got, err, callErr := callNative("manifestYamlFromJson", []interface{}{input})
	require.NoError(t, callErr)
	require.NoError(t, err)
	require.Equal(t, want, got)

	var decoded interface{}
	require.NoError(t, json.Unmarshal([]byte(input), &decoded))
	roundTrip, err, callErr := callNative("parseYaml", []interface{}{got})
	require.NoError(t, callErr)
	require.NoError(t, err)
	require.Equal(t, []interface{}{decoded}, roundTrip)
}
