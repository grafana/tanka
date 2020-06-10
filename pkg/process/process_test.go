package process

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"github.com/stretchr/testify/require"
)

func TestProcess(t *testing.T) {
	tests := []struct {
		name string
		spec v1alpha1.Spec

		deep interface{}
		flat manifest.List

		targets Matchers
		err     error
	}{
		{
			name: "regular",
			deep: testDataRegular().Deep,
			flat: mapToList(testDataRegular().Flat),
		},
		{
			name: "injectLabels",
			deep: testDataRegular().Deep,
			flat: mapToList(testDataRegular().Flat),
			spec: v1alpha1.Spec{
				InjectLabels: true,
			},
		},
		{
			name: "targets",
			deep: testDataDeep().Deep,
			flat: manifest.List{
				testDataDeep().Flat[".app.web.backend.server.grafana.deployment"],
				testDataDeep().Flat[".app.web.frontend.nodejs.express.service"],
			},
			targets: MustStrExps(
				`deployment/grafana`,
				`service/frontend`,
			),
		},
		{
			name: "targets-regex",
			deep: testDataDeep().Deep,
			flat: manifest.List{
				testDataDeep().Flat[".app.web.backend.server.grafana.deployment"],
				testDataDeep().Flat[".app.web.frontend.nodejs.express.deployment"],
			},
			targets: MustStrExps(`deployment/.*`),
		},
		{
			name: "targets-caseInsensitive",
			deep: testDataDeep().Deep,
			flat: manifest.List{
				testDataDeep().Flat[".app.web.backend.server.grafana.deployment"],
			},
			targets: MustStrExps(
				`DePlOyMeNt/GrAfAnA`,
			),
		},
		{
			name: "force-namespace",
			spec: v1alpha1.Spec{Namespace: "tanka"},
			deep: testDataFlat().Deep,
			flat: func() manifest.List {
				f := testDataFlat().Flat["."]
				f.Metadata()["namespace"] = "tanka"
				return manifest.List{f}
			}(),
		},
		{
			name: "custom-namespace",
			spec: v1alpha1.Spec{Namespace: "tanka"},
			deep: func() map[string]interface{} {
				d := testDataFlat().Deep.(map[string]interface{})
				d["metadata"].(map[string]interface{})["namespace"] = "custom"
				return d
			}(),
			flat: func() manifest.List {
				f := testDataFlat().Flat["."]
				f.Metadata()["namespace"] = "custom"
				return manifest.List{f}
			}(),
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			config := v1alpha1.New()
			config.Metadata.Name = "testdata"
			config.Spec = c.spec

			if config.Spec.InjectLabels {
				for i, m := range c.flat {
					m.Metadata().Labels()[LabelEnvironment] = config.Metadata.NameLabel()
					c.flat[i] = m
				}
			}

			got, err := Process(c.deep.(map[string]interface{}), *config, c.targets)
			require.Equal(t, c.err, err)

			Sort(c.flat)
			Sort(got)
			if diff := cmp.Diff(c.flat, got); diff != "" {
				t.Errorf("Process() mismatch:\n%s", diff)
			}
		})
	}
}

func mapToList(ms map[string]manifest.Manifest) manifest.List {
	l := make(manifest.List, 0, len(ms))
	for _, m := range ms {
		l = append(l, m)
	}
	return l
}

func TestProcessOrder(t *testing.T) {
	got := make([]manifest.List, 10)
	for i := 0; i < 10; i++ {
		r, err := Process(testDataDeep().Deep.(map[string]interface{}), *v1alpha1.New(), nil)
		require.NoError(t, err)
		got[i] = r
	}

	for i := 1; i < 10; i++ {
		require.Equal(t, got[0], got[i])
	}
}
