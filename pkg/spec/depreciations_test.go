package spec

import (
	"testing"

	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDeprecated checks that deprecated fields are still respected, but can be
// overwritten by the newer format.
func TestDeprecated(t *testing.T) {
	data := []byte(`
	{
      "spec": {
        "namespace": "new"
      },
	  "server": "https://127.0.0.1",
	  "team": "cool",
      "namespace": "old"
	}`)

	got, err := Parse(data, "test")
	require.Equal(t, ErrDeprecated{
		{old: "server", new: "spec.apiServer"},
		{old: "team", new: "metadata.labels.team"},
	}, err)

	want := v1alpha1.New()
	want.Spec.APIServer = "https://127.0.0.1"
	want.Spec.Namespace = "new"
	want.Metadata.Labels["team"] = "cool"
	want.Metadata.Name = "test"

	assert.Equal(t, want, got)
}
