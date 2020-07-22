package tanka

import (
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEvalJsonnet(t *testing.T) {

	tableTest := []struct {
		baseDir  string
		expected interface{}
	}{
		{baseDir: "./testdata/cases/array/",
			expected: []interface{}{
				[]interface{}{
					map[string]interface{}{"testCase": "nestedArray[0][0]"},
					map[string]interface{}{"testCase": "nestedArray[0][1]"},
				},
				[]interface{}{
					map[string]interface{}{"testCase": "nestedArray[1][0]"},
					map[string]interface{}{"testCase": "nestedArray[1][1]"},
				},
			},
		},
		{baseDir: "./testdata/cases/object/",
			expected: map[string]interface{}{
				"testCase": "object",
			},
		},
	}

	for _, test := range tableTest {
		m := make(map[string]string)
		data, e := evalJsonnet(test.baseDir, v1alpha1.New(), m)
		assert.NoError(t, e)
		assert.Equal(t, test.expected, data)
	}
}
