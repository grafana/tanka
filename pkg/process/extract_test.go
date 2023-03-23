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
		errMessage: `recursion ended on key "note" of type string which does not belong to a valid Kubernetes object
instead, it an attribute of the following object:

note: invalid because apiVersion and kind are missing


this object is not a valid Kubernetes object because: missing attribute "apiVersion"
`,
	},
	{
		name: "missing kind",
		data: testMissingAttribute(),
		errMessage: `recursion ended on key "apiVersion" of type string which does not belong to a valid Kubernetes object
instead, it an attribute of the following object:

apiVersion: v1
spec:
    ports:
        - port: 80
          protocol: TCP
          targetPort: 8080
    selector:
        app: deep


this object is not a valid Kubernetes object because: missing attribute "kind"
`,
	},
	{
		name: "bad kind",
		data: testBadKindType(),
		errMessage: `recursion ended on key "apiVersion" of type string which does not belong to a valid Kubernetes object
instead, it an attribute of the following object:

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


this object is not a valid Kubernetes object because: attribute "kind" is not a string, it is a float64
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
			} else {
				require.NoError(t, err)
				assert.EqualValues(t, c.data.Flat, extracted)
			}
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
