package tanka

import (
	"testing"

	"github.com/grafana/tanka/pkg/jsonnet"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestEvalJsonnet(t *testing.T) {
	cases := []struct {
		baseDir  string
		expected interface{}
	}{
		{
			baseDir: "./testdata/cases/array/",
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
		{
			baseDir: "./testdata/cases/object/",
			expected: map[string]interface{}{
				"testCase": "object",
			},
		},
	}

	for _, test := range cases {
		data, e := evalJsonnet(test.baseDir, v1alpha1.New(), jsonnet.Opts{})
		assert.NoError(t, e)
		assert.Equal(t, test.expected, data)
	}
}

func TestEval(t *testing.T) {
	cases := []struct {
		baseDir       string
		expectedError string
		expected      interface{}
	}{
		{
			baseDir:       "./testdata/cases/env-mismatch",
			expectedError: `reading spec.json: invalid metadata.name, must match generated "cases/env-mismatch"`,
		},
		{
			baseDir: "./testdata/cases/env",
			expected: map[string]interface{}{
				"tkName": "cases/env",
			},
		},
	}

	for _, test := range cases {
		raw, _, e := eval(test.baseDir, jsonnet.Opts{})
		if test.expectedError != "" {
			assert.Equal(t, test.expectedError, e.Error())
		} else {
			assert.NoError(t, e)
			assert.Equal(t, test.expected, raw)
		}
	}
}
