package process

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

func TestNamespace(t *testing.T) {
	cases := []struct {
		name          string
		namespace     string
		before, after manifest.Manifest
	}{
		// resource without a namespace: set it
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
			result := Namespace(manifest.List{c.before}, c.namespace)

			if diff := cmp.Diff(manifest.List{c.after}, result); diff != "" {
				t.Error(diff)
			}
		})
	}
}
