package tanka

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

func TestResourcesOutsideNamespace(t *testing.T) {
	resources := manifest.List{
		{"apiVersion": "v1", "kind": "Secret", "metadata": map[string]interface{}{"name": "inside", "namespace": "testing"}},
		{"apiVersion": "v1", "kind": "Secret", "metadata": map[string]interface{}{"name": "replicated", "namespace": "default"}},
		{"apiVersion": "v1", "kind": "Namespace", "metadata": map[string]interface{}{"name": "cluster-scoped"}},
		{"apiVersion": "v1", "kind": "ConfigMap", "metadata": map[string]interface{}{"name": "copy", "namespace": "review-app"}},
	}

	assert.Equal(t, []string{"default/Secret/replicated", "review-app/ConfigMap/copy"}, resourcesOutsideNamespace(resources, "testing"))
	assert.Empty(t, resourcesOutsideNamespace(resources[:1], "testing"))
}
