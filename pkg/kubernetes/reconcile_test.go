package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtract(t *testing.T) {
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
		{
			name: "array",
			data: testDataArray(),
		},
	}

	for _, c := range tests {
		t.Run(c.name, func(t *testing.T) {
			extracted, err := extract(c.data.deep)

			require.Equal(t, c.err, err)
			assert.EqualValues(t, c.data.flat, extracted)
		})
	}
}
