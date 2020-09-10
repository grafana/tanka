package process

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

func TestResourceDefaults(t *testing.T) {
	cases := []struct {
		name                string
		beforeAnnotations   map[string]string
		beforeLabels        map[string]string
		specAnnotations     map[string]string
		specLabels          map[string]string
		expectedAnnotations map[string]string
		expectedLabels      map[string]string
	}{
		// resource without a namespace: set it
		{
			name: "No change",
		},
		{
			name:                "Add annotation",
			specAnnotations:     map[string]string{"a": "b"},
			expectedAnnotations: map[string]string{"a": "b"},
		},
		{
			name:           "Add Label",
			specLabels:     map[string]string{"a": "b"},
			expectedLabels: map[string]string{"a": "b"},
		},
		{
			name:                "Add leaves existing",
			beforeAnnotations:   map[string]string{"1": "2"},
			beforeLabels:        map[string]string{"1": "2"},
			specAnnotations:     map[string]string{"a": "b"},
			specLabels:          map[string]string{"a": "b"},
			expectedAnnotations: map[string]string{"a": "b", "1": "2"},
			expectedLabels:      map[string]string{"a": "b", "1": "2"},
		},
		{
			name:                "Existing overrides spec",
			beforeAnnotations:   map[string]string{"a": "c", "1": "2"},
			beforeLabels:        map[string]string{"a": "c", "1": "2"},
			specAnnotations:     map[string]string{"a": "b"},
			specLabels:          map[string]string{"a": "b"},
			expectedAnnotations: map[string]string{"a": "c", "1": "2"},
			expectedLabels:      map[string]string{"a": "c", "1": "2"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			cfg := v1alpha1.Config{
				Spec: v1alpha1.Spec{
					ResourceDefaults: v1alpha1.ResourceDefaults{
						Annotations: c.specAnnotations,
						Labels:      c.specLabels,
					},
				},
			}

			before := manifest.Manifest{
				"kind": "Deployment",
			}
			for k, v := range c.beforeAnnotations {
				before.Metadata().Annotations()[k] = v
			}
			for k, v := range c.beforeLabels {
				before.Metadata().Labels()[k] = v
			}

			expected := manifest.Metadata{}
			for k, v := range c.expectedAnnotations {
				expected.Annotations()[k] = v
			}
			for k, v := range c.expectedLabels {
				expected.Labels()[k] = v
			}

			result := ResourceDefaults(manifest.List{before}, cfg)
			actual := result[0]
			if diff := cmp.Diff(expected, actual.Metadata()); diff != "" {
				t.Error(diff)
			}
		})
	}
}
