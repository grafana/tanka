package helm

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

func TestListAsMap(t *testing.T) {
	cases := []struct {
		name       string
		list       manifest.List
		result     map[string]interface{}
		nameFormat string
		err        error
	}{
		{
			name:       "simple",
			nameFormat: "", // test it properly defaults to DefaultNameFormat
			list: manifest.List{
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
			list: manifest.List{
				secret("foo", map[string]interface{}{"id": 1}),
				secret("foo", map[string]interface{}{"id": 2}),
			},
			err:    ErrorDuplicateName{name: "secret_foo", format: DefaultNameFormat},
			result: nil, // expect no result
		},
		{
			name:       "duplicate-custom",
			nameFormat: `{{ print .metadata.name "_" .data.id }}`,
			list: manifest.List{
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
			result, err := listAsMap(c.list, c.nameFormat)
			if err != c.err {
				t.Fatalf("err mismatch: want '%s' but got '%s'", c.err, err)
			}

			if diff := cmp.Diff(c.result, result); diff != "" {
				t.Fatal(diff)
			}
		})
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
