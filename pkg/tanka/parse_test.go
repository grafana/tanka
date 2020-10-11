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
		env      *v1alpha1.Config
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
			env: nil,
		},
		{
			baseDir: "./testdata/cases/object/",
			expected: map[string]interface{}{
				"testCase": "object",
			},
			env: nil,
		},
		{
			baseDir: "./testdata/cases/withspecjson/",
			expected: map[string]interface{}{
				"testCase": "object",
			},
			env: &v1alpha1.Config{
				APIVersion: v1alpha1.New().APIVersion,
				Kind:       v1alpha1.New().Kind,
				Metadata: v1alpha1.Metadata{
					Name:   "cases/withspecjson",
					Labels: v1alpha1.New().Metadata.Labels,
				},
				Spec: v1alpha1.Spec{
					APIServer: "https://localhost",
					Namespace: "withspec",
				},
				Data: map[string]interface{}{
					"testCase": "object",
				},
			},
		},
		{
			baseDir: "./testdata/cases/withspecjson/main.jsonnet",
			expected: map[string]interface{}{
				"testCase": "object",
			},
			env: &v1alpha1.Config{
				APIVersion: v1alpha1.New().APIVersion,
				Kind:       v1alpha1.New().Kind,
				Metadata: v1alpha1.Metadata{
					Name:   "cases/withspecjson",
					Labels: v1alpha1.New().Metadata.Labels,
				},
				Spec: v1alpha1.Spec{
					APIServer: "https://localhost",
					Namespace: "withspec",
				},
				Data: map[string]interface{}{
					"testCase": "object",
				},
			},
		},
		{
			baseDir: "./testdata/cases/withenv/main.jsonnet",
			expected: v1alpha1.Config{
				APIVersion: v1alpha1.New().APIVersion,
				Kind:       v1alpha1.New().Kind,
				Metadata: v1alpha1.Metadata{
					Name: "withenv",
				},
				Spec: v1alpha1.Spec{
					APIServer: "https://localhost",
					Namespace: "withenv",
				},
				Data: map[string]interface{}{
					"testCase": "object",
				},
			},
			env: &v1alpha1.Config{
				APIVersion: v1alpha1.New().APIVersion,
				Kind:       v1alpha1.New().Kind,
				Metadata: v1alpha1.Metadata{
					Name: "withenv",
				},
				Spec: v1alpha1.Spec{
					APIServer: "https://localhost",
					Namespace: "withenv",
				},
				Data: map[string]interface{}{
					"testCase": "object",
				},
			},
		},
	}

	for _, test := range cases {
		data, env, e := evalJsonnet(test.baseDir, jsonnet.Opts{})
		assert.NoError(t, e)
		assert.Equal(t, test.expected, data)
		assert.Equal(t, test.env, env)
	}
}
