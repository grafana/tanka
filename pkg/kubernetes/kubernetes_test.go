package kubernetes

import (
	"testing"

	"github.com/grafana/tanka/pkg/config/v1alpha1"
	"github.com/stretchr/objx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReconcile(t *testing.T) {
	tests := []struct {
		name string
		k    *Kubernetes

		data    testData
		targets []string
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
			targets: []string{"deployment/nginx", "service/frontend"},
		},
		{
			name: "force-namespace",
			k:    &Kubernetes{Spec: v1alpha1.Spec{Namespace: "tanka"}},
			data: testData{
				deep: testDataFlat().deep,
				flat: func() []map[string]interface{} {
					f := objx.New(testDataFlat().flat.([]map[string]interface{})[0])
					f.Set("metadata.namespace", "tanka")
					return []map[string]interface{}{f}
				}(),
			},
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			got, err := c.k.Reconcile(c.data.deep.(map[string]interface{}), c.targets...)

			require.Equal(t, c.err, err)

			flat := c.data.flat.([]map[string]interface{})
			assert.Equal(t, msisToManifests(flat), got)
		})
	}
}

func TestReconcileOrder(t *testing.T) {
	got := make([][]Manifest, 10)
	k := &Kubernetes{}
	for i := 0; i < 10; i++ {
		r, err := k.Reconcile(testDataDeep().deep.(map[string]interface{}))
		require.NoError(t, err)
		got[i] = r
	}

	for i := 1; i < 10; i++ {
		require.Equal(t, got[0], got[i])
	}
}

func msisToManifests(msis []map[string]interface{}) []Manifest {
	ms := make([]Manifest, len(msis))
	for i, msi := range msis {
		ms[i] = Manifest(msi)
	}
	return ms
}
