package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var extractTestCases = []struct {
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
		err:  ErrorPrimitiveReached{path: ".service", key: "note", primitive: "invalid because apiVersion and kind are missing"},
	},
	{
		name: "deep",
		data: testDataDeep(),
	},
	{
		name: "array",
		data: testDataArray(),
	},
	{
		name: "nil",
		data: func() testData {
			d := testDataRegular()
			d.Deep.(map[string]interface{})["disabledObject"] = nil
			return d
		}(),
		err: nil, // we expect no error, just the result of testDataRegular
	},
}

func TestExtract(t *testing.T) {
	for _, c := range extractTestCases {
		t.Run(c.name, func(t *testing.T) {
			extracted, err := Extract(c.data.Deep)

			require.Equal(t, c.err, err)
			assert.EqualValues(t, c.data.Flat, extracted)
		})
	}
}

func BenchmarkExtract(b *testing.B) {
	for _, c := range extractTestCases {
		if c.err != nil {
			continue
		}
		b.Run(c.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				// nolint:errcheck
				Extract(c.data.Deep)
			}
		})
	}
}
