package kubernetes

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/kubernetes/util"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

func TestReconcile(t *testing.T) {
	tests := []struct {
		name string
		spec v1alpha1.Spec

		deep interface{}
		flat manifest.List

		targets []*regexp.Regexp
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
			targets: util.MustCompileTargetExps(
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
			targets: util.MustCompileTargetExps(`deployment/.*`),
		},
		{
			name: "targets-caseInsensitive",
			deep: testDataDeep().Deep,
			flat: manifest.List{
				testDataDeep().Flat[".app.web.backend.server.grafana.deployment"],
			},
			targets: util.MustCompileTargetExps(
				`DePlOyMeNt/GrAfAnA`,
			),
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

			got, err := Reconcile(c.deep.(map[string]interface{}), *config, c.targets)
			require.Equal(t, c.err, err)

			assert.ElementsMatch(t, c.flat, got)
		})
	}
}

func TestReconcileOrder(t *testing.T) {
	got := make([]manifest.List, 10)
	for i := 0; i < 10; i++ {
		r, err := Reconcile(testDataDeep().Deep.(map[string]interface{}), *v1alpha1.New(), nil)
		require.NoError(t, err)
		got[i] = r
	}

	for i := 1; i < 10; i++ {
		require.Equal(t, got[0], got[i])
	}
}

func mapToList(ms map[string]manifest.Manifest) manifest.List {
	l := make(manifest.List, 0, len(ms))
	for _, m := range ms {
		l = append(l, m)
	}
	return l
}
