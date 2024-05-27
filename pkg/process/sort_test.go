package process

import (
	"testing"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/stretchr/testify/require"
)

func TestSort(t *testing.T) {
	cases := []struct {
		raw   manifest.List
		state manifest.List
		err   error
	}{
		{
			// sorting by kinds in the `kindOrder` list
			raw: manifest.List{
				mkobj("Service", "service", "default"),
				mkobj("Deployment", "deployment", "default"),
				mkobj("CustomResourceDefinition", "crd", ""),
			},
			state: manifest.List{
				mkobj("CustomResourceDefinition", "crd", ""),
				mkobj("Service", "service", "default"),
				mkobj("Deployment", "deployment", "default"),
			},
		},
		{
			// alphabtical sorting by kinds outside `kindOrder` list
			raw: manifest.List{
				mkobj("B", "b", "default"),
				mkobj("C", "c", "default"),
				mkobj("A", "a", "default"),
			},
			state: manifest.List{
				mkobj("A", "a", "default"),
				mkobj("B", "b", "default"),
				mkobj("C", "c", "default"),
			},
		},
		{
			// sorting by the namespace if kinds match
			raw: manifest.List{
				mkobj("Service", "service", "default2"),
				mkobj("Service", "service", "default"),
				mkobj("Service", "service", "default1"),
			},
			state: manifest.List{
				mkobj("Service", "service", "default"),
				mkobj("Service", "service", "default1"),
				mkobj("Service", "service", "default2"),
			},
		},
		{
			// sorting by the names if both kinds and namespaces match
			raw: manifest.List{
				mkobj("Service", "service2", "default"),
				mkobj("Service", "service", "default"),
				mkobj("Service", "service1", "default"),
			},
			state: manifest.List{
				mkobj("Service", "service", "default"),
				mkobj("Service", "service1", "default"),
				mkobj("Service", "service2", "default"),
			},
		},
		{
			// sorting by the names if both kinds match and there are no namespaces
			raw: manifest.List{
				mkobj("CustomResourceDefinition", "crd2", ""),
				mkobj("CustomResourceDefinition", "crd", ""),
				mkobj("CustomResourceDefinition", "crd1", ""),
			},
			state: manifest.List{
				mkobj("CustomResourceDefinition", "crd", ""),
				mkobj("CustomResourceDefinition", "crd1", ""),
				mkobj("CustomResourceDefinition", "crd2", ""),
			},
		},
		{
			// sorting by the generate name prefix if everything else is the same
			raw: manifest.List{
				mkGenerateObj("CustomResourceDefinition", "crd2-", ""),
				mkGenerateObj("CustomResourceDefinition", "crd-", ""),
				mkGenerateObj("CustomResourceDefinition", "crd1-", ""),
			},
			state: manifest.List{
				mkGenerateObj("CustomResourceDefinition", "crd-", ""),
				mkGenerateObj("CustomResourceDefinition", "crd1-", ""),
				mkGenerateObj("CustomResourceDefinition", "crd2-", ""),
			},
		},
		{
			raw: manifest.List{
				mkobj("Deployment", "b", "a"),
				mkobj("ConfigMap", "a", "a"),
				mkobj("Issuer", "a", "a"),
				mkobj("Service", "b", "a"),
				mkobj("Service", "a", "a"),
				mkobj("Deployment", "a", "a"),
				mkobj("ConfigMap", "a", "b"),
				mkobj("Issuer", "b", "a"),
				mkobj("Service", "b", "b"),
				mkobj("Deployment", "a", "b"),
				mkobj("ConfigMap", "b", "a"),
				mkobj("Issuer", "a", "b"),
			},
			state: manifest.List{
				mkobj("ConfigMap", "a", "a"),
				mkobj("ConfigMap", "b", "a"),
				mkobj("ConfigMap", "a", "b"),
				mkobj("Service", "a", "a"),
				mkobj("Service", "b", "a"),
				mkobj("Service", "b", "b"),
				mkobj("Deployment", "a", "a"),
				mkobj("Deployment", "b", "a"),
				mkobj("Deployment", "a", "b"),
				mkobj("Issuer", "a", "a"),
				mkobj("Issuer", "b", "a"),
				mkobj("Issuer", "a", "b"),
			},
		},
	}

	for _, c := range cases {
		Sort(c.raw)
		require.Equal(t, c.state, c.raw)
	}
}

func mkobj(kind string, name string, ns string) map[string]interface{} {
	ret := map[string]interface{}{
		"kind":       kind,
		"apiVersion": "apiversion",
		"metadata": map[string]interface{}{
			"name": name,
		},
	}
	if ns != "" {
		ret["metadata"].(map[string]interface{})["namespace"] = ns
	}

	return ret
}

func mkGenerateObj(kind string, generateName string, ns string) map[string]interface{} {
	result := mkobj(kind, "", ns)
	result["metadata"].(map[string]interface{})["generateName"] = generateName
	return result
}
