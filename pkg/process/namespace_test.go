package process

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/kubernetes/resources"
)

func TestNamespace(t *testing.T) {
	cases := []struct {
		name          string
		namespace     string
		before, after manifest.Manifest
	}{
		// namespaced resource without a namespace: set it
		{
			name:      "simple/namespaced",
			namespace: "testing",
			before: manifest.Manifest{
				"kind": "Deployment",
			},
			after: manifest.Manifest{
				"kind": "Deployment",
				"metadata": map[string]interface{}{
					"namespace": "testing",
				},
			},
		},

		// non-namespaced resource without a namespace: not set it
		{
			name:      "simple/cluster-wide",
			namespace: "testing",
			before: manifest.Manifest{
				"kind":     "ClusterRole",
				"metadata": map[string]interface{}{},
			},
			after: manifest.Manifest{
				"kind":     "ClusterRole",
				"metadata": map[string]interface{}{},
			},
		},

		// resource with a namespace: ignore it
		{
			name:      "already-present",
			namespace: "ignored",
			before: manifest.Manifest{
				"kind":     "Deployment",
				"metadata": map[string]interface{}{"namespace": "mycoolnamespace"},
			},
			after: manifest.Manifest{
				"kind":     "Deployment",
				"metadata": map[string]interface{}{"namespace": "mycoolnamespace"},
			},
		},

		// custom resource: do nothing
		{
			name:      "custom-resource",
			namespace: "testing",
			before: manifest.Manifest{
				"kind":     "MyCoolThing",
				"metadata": map[string]interface{}{},
			},
			after: manifest.Manifest{
				"kind":     "MyCoolThing",
				"metadata": map[string]interface{}{},
			},
		},

		// custom resource, explicit set
		{
			name:      "custom-resource-explicit",
			namespace: "testing",
			before: manifest.Manifest{
				"kind": "MyCoolThing",
				"metadata": map[string]interface{}{
					"annotations": map[string]interface{}{
						"tanka.dev/namespaced": "true",
					},
				},
			},
			after: manifest.Manifest{
				"kind": "MyCoolThing",
				"metadata": map[string]interface{}{
					"annotations": map[string]string{
						"tanka.dev/namespaced": "true",
					},
					"namespace": "testing",
				},
			},
		},

		// empty default ns: do nothing
		{
			name:      "no-default",
			namespace: "",
			before: manifest.Manifest{
				"kind":     "Deployment",
				"metadata": map[string]interface{}{},
			},
			after: manifest.Manifest{
				"kind":     "Deployment",
				"metadata": map[string]interface{}{},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			result := Namespace(manifest.List{c.before}, c.namespace, resources.StaticStore)

			if diff := cmp.Diff(manifest.List{c.after}, result); diff != "" {
				t.Errorf("Namespace() mismatch:\n%s", diff)
			}
		})
	}
}
