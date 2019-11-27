package kubernetes

import (
	"regexp"
	"testing"

	"github.com/stretchr/objx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
)

func TestPrepare(t *testing.T) {
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
			deep: testDataRegular().deep,
			flat: mapToList(testDataRegular().flat),
		},
		{
			name: "targets",
			deep: testDataDeep().deep,
			flat: manifest.List{
				testDataDeep().flat[".app.web.backend.server.nginx.deployment"],
				testDataDeep().flat[".app.web.frontend.nodejs.express.service"],
			},
			targets: []*regexp.Regexp{
				regexp.MustCompile("deployment/nginx"),
				regexp.MustCompile("service/frontend"),
			},
		},
		{
			name: "targets-regex",
			deep: testDataDeep().deep,
			flat: manifest.List{
				testDataDeep().flat[".app.web.backend.server.nginx.deployment"],
				testDataDeep().flat[".app.web.frontend.nodejs.express.deployment"],
			},
			targets: []*regexp.Regexp{regexp.MustCompile("deployment/.*")},
		},
		{
			name: "force-namespace",
			spec: v1alpha1.Spec{Namespace: "tanka"},
			deep: testDataFlat().deep,
			flat: func() manifest.List {
				f := testDataFlat().flat["."]
				f.Metadata()["namespace"] = "tanka"
				return manifest.List{f}
			}(),
		},
		{
			name: "custom-namespace",
			spec: v1alpha1.Spec{Namespace: "tanka"},
			deep: func() map[string]interface{} {
				d := objx.New(testDataFlat().deep)
				d.Set("metadata.namespace", "custom")
				return d
			}(),
			flat: func() manifest.List {
				f := testDataFlat().flat["."]
				f.Metadata()["namespace"] = "custom"
				return manifest.List{f}
			}(),
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			got, err := Prepare(c.deep.(map[string]interface{}), c.spec, c.targets)

			require.Equal(t, c.err, err)
			assert.ElementsMatch(t, c.flat, got)
		})
	}
}

func TestPrepareOrder(t *testing.T) {
	got := make([]manifest.List, 10)
	for i := 0; i < 10; i++ {
		r, err := Prepare(testDataDeep().deep.(map[string]interface{}), v1alpha1.Spec{}, nil)
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
