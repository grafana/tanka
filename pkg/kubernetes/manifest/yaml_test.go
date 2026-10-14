package manifest

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func TestYAMLSerialization(t *testing.T) {
	configMap := Manifest{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata":   map[string]interface{}{"name": "example"},
		"data": map[string]interface{}{
			"boolean": "true",
			"number":  "001",
			"script":  "echo first\necho second\n",
		},
	}
	deployment := Manifest{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata":   map[string]interface{}{"name": "example"},
		"spec": map[string]interface{}{
			"replicas": float64(3),
			"paused":   false,
			"template": map[string]interface{}{
				"spec": map[string]interface{}{
					"containers": []interface{}{
						map[string]interface{}{"name": "app", "image": "nginx:1.14.2"},
					},
				},
			},
		},
	}
	const configMapYAML = `apiVersion: v1
data:
  boolean: "true"
  number: "001"
  script: |
    echo first
    echo second
kind: ConfigMap
metadata:
  name: example
`
	const deploymentYAML = `apiVersion: apps/v1
kind: Deployment
metadata:
  name: example
spec:
  paused: false
  replicas: 3
  template:
    spec:
      containers:
      - image: nginx:1.14.2
        name: app
`

	require.Equal(t, configMapYAML, configMap.String())
	require.Equal(t, deploymentYAML, deployment.String())

	list := List{configMap, deployment}
	stream := list.String()
	require.Equal(t, configMapYAML+"---\n"+deploymentYAML, stream)

	decoder := yaml.NewDecoder(strings.NewReader(stream))
	for _, want := range list {
		var got Manifest
		require.NoError(t, decoder.Decode(&got))
		require.Equal(t, want, got)
	}
	var extra Manifest
	require.ErrorIs(t, decoder.Decode(&extra), io.EOF)
}
