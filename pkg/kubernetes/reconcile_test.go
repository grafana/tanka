package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWalkJSON(t *testing.T) {
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
			err:  ErrorPrimitiveReached{path: ".nginx.service", key: "note", primitive: "invalid because apiVersion and kind are missing"},
		},
		{
			name: "deep",
			data: testDataDeep(),
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			got, err := walkJSON(c.data.deep.(map[string]interface{}), "")
			require.Equal(t, c.err, err)
			assert.ElementsMatch(t, c.data.flat, got)
		})
	}
}
