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

func TestReconcile(t *testing.T) {
	tests := []struct {
		name string
		spec v1alpha1.Spec

		data    testData
		targets []*regexp.Regexp
		err     error
	}{
		{
			name: "regular",
			data: testDataRegular(),
		},
		{
			name: "targets",
			data: testData{
				deep: testDataDeep().deep,
				flat: []map[string]interface{}{
					testDataDeep().flat.([]map[string]interface{})[0], // deployment/nginx
					testDataDeep().flat.([]map[string]interface{})[1], // service/frontend
				},
			},
			targets: []*regexp.Regexp{
				regexp.MustCompile("deployment/nginx"),
				regexp.MustCompile("service/frontend"),
			},
		},
		{
			name: "targets-regex",
			data: testData{
				deep: testDataDeep().deep,
				flat: []map[string]interface{}{
					testDataDeep().flat.([]map[string]interface{})[2], // deployment/frontend
					testDataDeep().flat.([]map[string]interface{})[0], // deployment/nginx
				},
			},
			targets: []*regexp.Regexp{regexp.MustCompile("deployment/.*")},
		},
		{
			name: "force-namespace",
			spec: v1alpha1.Spec{Namespace: "tanka"},
			data: testData{
				deep: testDataFlat().deep,
				flat: func() []map[string]interface{} {
					f := objx.New(testDataFlat().flat.([]map[string]interface{})[0])
					f.Set("metadata.namespace", "tanka")
					return []map[string]interface{}{f}
				}(),
			},
		},
		{
			name: "custom-namespace",
			spec: v1alpha1.Spec{Namespace: "tanka"},
			data: testData{
				deep: func() map[string]interface{} {
					d := objx.New(testDataFlat().deep)
					d.Set("metadata.namespace", "custom")
					return d
				}(),
				flat: func() []map[string]interface{} {
					f := objx.New(testDataFlat().flat.([]map[string]interface{})[0])
					f.Set("metadata.namespace", "custom")
					return []map[string]interface{}{f}
				}(),
			},
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			got, err := Reconcile(c.data.deep.(map[string]interface{}), c.spec, c.targets)

			require.Equal(t, c.err, err)

			flat := c.data.flat.([]map[string]interface{})
			assert.Equal(t, msisToManifests(flat), got)
		})
	}
}

func TestReconcileOrder(t *testing.T) {
	got := make([]manifest.List, 10)
	for i := 0; i < 10; i++ {
		r, err := Reconcile(testDataDeep().deep.(map[string]interface{}), v1alpha1.Spec{}, nil)
		require.NoError(t, err)
		got[i] = r
	}

	for i := 1; i < 10; i++ {
		require.Equal(t, got[0], got[i])
	}
}

func msisToManifests(msis []map[string]interface{}) manifest.List {
	ms := make(manifest.List, len(msis))
	for i, msi := range msis {
		ms[i] = manifest.Manifest(msi)
	}
	return ms
}
