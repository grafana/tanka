package manifest

import (
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

// UnmarshalExpect defines the expected Unmarshal result. Types are very
// important here, only nil, float64, bool, string, map[string]interface{} and
// []interface{} may exist.
var UnmarshalExpect = Manifest{
	"apiVersion": string("apps/v1"),
	"kind":       string("Deployment"),
	"metadata": map[string]interface{}{
		"name": string("MyDeployment"),
	},
	"spec": map[string]interface{}{
		"replicas": float64(3),
		"template": map[string]interface{}{
			"spec": map[string]interface{}{
				"containers": []interface{}{
					map[string]interface{}{
						"name":  string("nginx"),
						"image": string("nginx:1.14.2"),
					},
				},
			},
		},
	},
}

func TestUnmarshalJSON(t *testing.T) {
	const data = `
{
   "apiVersion":"apps/v1",
   "kind":"Deployment",
   "metadata":{
	  "name":"MyDeployment"
   },
   "spec":{
	  "replicas":3,
	  "template":{
		 "spec":{
			"containers":[
			   {
				  "name":"nginx",
				  "image":"nginx:1.14.2"
			   }
			]
		 }
	  }
   }
}
`

	var m Manifest
	err := json.Unmarshal([]byte(data), &m)
	require.NoError(t, err)

	if s := cmp.Diff(UnmarshalExpect, m); s != "" {
		t.Error(s)
	}
}

func TestUnmarshalYAML(t *testing.T) {
	const data = `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: MyDeployment
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
`

	var m Manifest
	err := yaml.Unmarshal([]byte(data), &m)
	require.NoError(t, err)

	if s := cmp.Diff(UnmarshalExpect, m); s != "" {
		t.Error(s)
	}
}

func TestListAsMap(t *testing.T) {
	cases := []struct {
		name       string
		list       List
		result     map[string]interface{}
		nameFormat string
		err        error
	}{
		{
			name:       "simple",
			nameFormat: "", // test it properly defaults to DefaultNameFormat
			list: List{
				configMap("foo"),
				deployment("bar"),
			},
			result: map[string]interface{}{
				"config_map_foo": configMap("foo"),
				"deployment_bar": deployment("bar"),
			},
		},
		{
			name: "duplicate-default",
			list: List{
				secret("foo", map[string]interface{}{"id": 1}),
				secret("foo", map[string]interface{}{"id": 2}),
			},
			err:    ErrorDuplicateName{name: "secret_foo", format: DefaultNameFormat},
			result: nil, // expect no result
		},
		{
			name:       "duplicate-custom",
			nameFormat: `{{ print .metadata.name "_" .data.id }}`,
			list: List{
				secret("foo", map[string]interface{}{"id": 1}),
				secret("foo", map[string]interface{}{"id": 2}),
			},
			result: map[string]interface{}{
				"foo_1": secret("foo", map[string]interface{}{"id": 1}),
				"foo_2": secret("foo", map[string]interface{}{"id": 2}),
			},
			err: nil, // expect no error
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result, err := ListAsMap(c.list, c.nameFormat)
			if err != c.err {
				t.Fatalf("err mismatch: want '%s' but got '%s'", c.err, err)
			}

			if diff := cmp.Diff(c.result, result); diff != "" {
				t.Fatal(diff)
			}
		})
	}
}

func TestManifestMarshalMultiline(t *testing.T) {
	const data = `
{
   "apiVersion":"core/v1",
   "kind":"ConfigMap",
   "metadata":{
	  "name":"MyConfigMap"
   },
   "data":{
	"script.sh": "#/bin/sh\nset -e\n\n# This is a sample secript as configmap\n\necho \"test\"if test -f 'test'; then\n\techo \"If test\"\nfi"
   }
}
`

	const expectedYAML = `apiVersion: core/v1
data:
    script.sh: |-
        #/bin/sh
        set -e

        # This is a sample secript as configmap

        echo "test"if test -f 'test'; then
        	echo "If test"
        fi
kind: ConfigMap
metadata:
    name: MyConfigMap
`
	var m Manifest
	err := json.Unmarshal([]byte(data), &m)
	require.NoError(t, err)

	outYAML := m.String()
	if diff := cmp.Diff(outYAML, expectedYAML); diff != "" {
		t.Fatal(diff)
	}
}

func configMap(name string) map[string]interface{} {
	return map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ConfigMap",
		"metadata":   map[string]interface{}{"name": name},
		"data":       map[string]interface{}{},
	}
}

func deployment(name string) map[string]interface{} {
	return map[string]interface{}{
		"apiVersion": "apps/v1",
		"kind":       "Deployment",
		"metadata":   map[string]interface{}{"name": name},
		"spec":       map[string]interface{}{},
	}
}

func secret(name string, data map[string]interface{}) map[string]interface{} {
	if data == nil {
		data = map[string]interface{}{}
	}

	return map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "Secret",
		"metadata":   map[string]interface{}{"name": name},
		"data":       data,
	}
}
