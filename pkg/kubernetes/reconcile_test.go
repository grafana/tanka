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
			err:  ErrorPrimitiveReached{path: "nginx/service/note", primitive: "invalid because apiVersion and kind are missing"},
		},
		{
			name: "deep",
			data: testDataDeep(),
		},
		{
			name: "array",
			data: testDataArray(),
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			manifests, err := walkJSON(c.data.deep)

			expectedManifests := []Manifest{}
			for _, manifest := range c.data.flat {
				expectedManifests = append(expectedManifests, manifest)
			}

			require.Equal(t, c.err, err)
			assert.ElementsMatch(t, expectedManifests, manifests)
		})
	}
}
