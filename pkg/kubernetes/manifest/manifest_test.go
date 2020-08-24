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
