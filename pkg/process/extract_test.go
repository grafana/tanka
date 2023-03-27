package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var extractTestCases = []struct {
	name       string
	data       testData
	errMessage string
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
		errMessage: `found invalid Kubernetes object (at .service): missing attribute "apiVersion"

note: invalid because apiVersion and kind are missing
`,
	},
	{
		name: "missing kind",
		data: testMissingAttribute(),
		errMessage: `found invalid Kubernetes object (at .service): missing attribute "kind"

apiVersion: v1
spec:
    ports:
        - port: 80
          protocol: TCP
          targetPort: 8080
    selector:
        app: deep
`,
	},
	{
		name: "bad kind",
		data: testBadKindType(),
		errMessage: `found invalid Kubernetes object (at .deployment): attribute "kind" is not a string, it is a float64

apiVersion: apps/v1
kind: 3000
metadata:
    name: grafana
spec:
    replicas: 1
    template:
        containers:
            - image: grafana/grafana
              name: grafana
        metadata:
            labels:
                app: grafana
`,
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
		errMessage: "", // we expect no error, just the result of testDataRegular
	},
}

func TestExtract(t *testing.T) {
	for _, c := range extractTestCases {
		t.Run(c.name, func(t *testing.T) {
			extracted, err := Extract(c.data.Deep)

			if c.errMessage != "" {
				require.Error(t, err)
				assert.Equal(t, c.errMessage, err.Error())
				return
			}

			require.NoError(t, err)
			assert.EqualValues(t, c.data.Flat, extracted)
		})
	}
}

func BenchmarkExtract(b *testing.B) {
	for _, c := range extractTestCases {
		if c.errMessage != "" {
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
