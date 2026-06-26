package process

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
)

func TestInjectLabels(t *testing.T) {
	cases := []struct {
		name           string
		beforeLabels   map[string]string
		injectLabels   map[string]string
		expectedLabels map[string]string
	}{
		{
			name: "No labels is a no-op",
		},
		{
			name:           "Add label",
			injectLabels:   map[string]string{"a": "b"},
			expectedLabels: map[string]string{"a": "b"},
		},
		{
			name:           "Add multiple labels",
			injectLabels:   map[string]string{"a": "b", "c": "d"},
			expectedLabels: map[string]string{"a": "b", "c": "d"},
		},
		{
			name:           "Add leaves unrelated labels",
			beforeLabels:   map[string]string{"1": "2"},
			injectLabels:   map[string]string{"a": "b"},
			expectedLabels: map[string]string{"a": "b", "1": "2"},
		},
		{
			// Unlike ResourceDefaults, injected labels override existing values.
			name:           "Injected labels override existing",
			beforeLabels:   map[string]string{"a": "c", "1": "2"},
			injectLabels:   map[string]string{"a": "b"},
			expectedLabels: map[string]string{"a": "b", "1": "2"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			before := manifest.Manifest{
				"kind": "Deployment",
			}
			for k, v := range c.beforeLabels {
				before.Metadata().Labels()[k] = v
			}

			expected := manifest.Metadata{}
			for k, v := range c.expectedLabels {
				expected.Labels()[k] = v
			}

			result := InjectLabels(manifest.List{before}, c.injectLabels)
			actual := result[0]
			if diff := cmp.Diff(expected, actual.Metadata()); diff != "" {
				t.Error(diff)
			}
		})
	}
}

func TestInjectLabelsAllManifests(t *testing.T) {
	list := manifest.List{
		{"kind": "Deployment"},
		{"kind": "Service"},
	}

	result := InjectLabels(list, map[string]string{"created-by": "tester"})

	for _, m := range result {
		if got := m.Metadata().Labels()["created-by"]; got != "tester" {
			t.Errorf("expected label on every manifest, %s has %v", m.Kind(), got)
		}
	}
}
