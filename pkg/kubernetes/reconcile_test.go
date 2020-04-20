package kubernetes

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

func TestExtract(t *testing.T) {
	tests := []struct {
		name string
		data testData
		err  error
	}{
		{
			name: "regular",
			data: testDataRegular(),
		},
		{
			name: "flat",
			data: testDataFlat(),
		},
		{
			name: "primitive",
			data: testDataPrimitive(),
			err:  ErrorPrimitiveReached{path: ".service", key: "note", primitive: "invalid because apiVersion and kind are missing"},
		},
		{
			name: "deep",
			data: testDataDeep(),
		},
		{
			name: "array",
			data: testDataArray(),
		},
		{
			name: "nil",
			data: func() testData {
				d := testDataRegular()
				d.Deep.(map[string]interface{})["disabledObject"] = nil
				return d
			}(),
			err: nil, // we expect no error, just the result of testDataRegular
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			extracted, err := extract(c.data.Deep)

			require.Equal(t, c.err, err)
			assert.EqualValues(t, c.data.Flat, extracted)
		})
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

func TestReconcileSorting(t *testing.T) {
	tests := []struct {
		raw     map[string]interface{}
		targets []*regexp.Regexp
		state   manifest.List
		err     error
	}{
		{
			// sorting by kinds in the `kindOrder` list
			raw: map[string]interface{}{
				"a": mkobj("Service", "service", "default"),
				"b": mkobj("Deployment", "deployment", "default"),
				"c": mkobj("CustomResourceDefinition", "crd", ""),
			},
			state: manifest.List{
				mkobj("CustomResourceDefinition", "crd", ""),
				mkobj("Service", "service", "default"),
				mkobj("Deployment", "deployment", "default"),
			},
		},
		{
			// alphabtical sorting by kinds outside `kindOrder` list
			raw: map[string]interface{}{
				"a": mkobj("B", "b", "default"),
				"b": mkobj("C", "c", "default"),
				"c": mkobj("A", "a", "default"),
			},
			state: manifest.List{
				mkobj("A", "a", "default"),
				mkobj("B", "b", "default"),
				mkobj("C", "c", "default"),
			},
		},
		{
			// sorting by the namespace if kinds match
			raw: map[string]interface{}{
				"a": mkobj("Service", "service", "default2"),
				"b": mkobj("Service", "service", "default"),
				"c": mkobj("Service", "service", "default1"),
			},
			state: manifest.List{
				mkobj("Service", "service", "default"),
				mkobj("Service", "service", "default1"),
				mkobj("Service", "service", "default2"),
			},
		},
		{
			// sorting by the names if both kinds and namespaces match
			raw: map[string]interface{}{
				"a": mkobj("Service", "service2", "default"),
				"b": mkobj("Service", "service", "default"),
				"c": mkobj("Service", "service1", "default"),
			},
			state: manifest.List{
				mkobj("Service", "service", "default"),
				mkobj("Service", "service1", "default"),
				mkobj("Service", "service2", "default"),
			},
		},
		{
			// sorting by the names if both kinds match and there are no namespaces
			raw: map[string]interface{}{
				"a": mkobj("CustomResourceDefinition", "crd2", ""),
				"b": mkobj("CustomResourceDefinition", "crd", ""),
				"c": mkobj("CustomResourceDefinition", "crd1", ""),
			},
			state: manifest.List{
				mkobj("CustomResourceDefinition", "crd", ""),
				mkobj("CustomResourceDefinition", "crd1", ""),
				mkobj("CustomResourceDefinition", "crd2", ""),
			},
		},
		{
			raw: map[string]interface{}{
				"a": mkobj("Deployment", "b", "a"),
				"b": mkobj("ConfigMap", "a", "a"),
				"c": mkobj("Issuer", "a", "a"),
				"d": mkobj("Service", "b", "a"),
				"e": mkobj("Service", "a", "a"),
				"f": mkobj("Deployment", "a", "a"),
				"g": mkobj("ConfigMap", "a", "b"),
				"h": mkobj("Issuer", "b", "a"),
				"i": mkobj("Service", "b", "b"),
				"j": mkobj("Deployment", "a", "b"),
				"k": mkobj("ConfigMap", "b", "a"),
				"l": mkobj("Issuer", "a", "b"),
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

	for _, test := range tests {
		res, err := Reconcile(test.raw, *v1alpha1.New(), test.targets)

		require.NoError(t, err)
		require.Equal(t, test.state, res)
	}
}
