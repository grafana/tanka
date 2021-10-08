package tanka

import (
	"testing"

	"github.com/grafana/tanka/pkg/kubernetes/manifest"
	"github.com/grafana/tanka/pkg/spec/v1alpha1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	cases := []struct {
		name     string
		baseDir  string
		expected interface{}
		env      *v1alpha1.Environment
	}{
		{
			name:    "static",
			baseDir: "./testdata/cases/withspecjson/",
			expected: manifest.List{{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]interface{}{
					"name":      "config",
					"namespace": "withspec",
				},
			}},
			env: &v1alpha1.Environment{
				APIVersion: v1alpha1.New().APIVersion,
				Kind:       v1alpha1.New().Kind,
				Metadata: v1alpha1.Metadata{
					Name:      "cases/withspecjson",
					Namespace: "cases/withspecjson/main.jsonnet",
					Labels:    v1alpha1.New().Metadata.Labels,
				},
				Spec: v1alpha1.Spec{
					APIServer: "https://localhost",
					Namespace: "withspec",
				},
				Data: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata":   map[string]interface{}{"name": "config", "namespace": "withspec"},
				},
			},
		},
		{
			name:    "static-filename",
			baseDir: "./testdata/cases/withspecjson/main.jsonnet",
			expected: manifest.List{{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]interface{}{
					"name":      "config",
					"namespace": "withspec",
				},
			}},
			env: &v1alpha1.Environment{
				APIVersion: v1alpha1.New().APIVersion,
				Kind:       v1alpha1.New().Kind,
				Metadata: v1alpha1.Metadata{
					Name:      "cases/withspecjson",
					Namespace: "cases/withspecjson/main.jsonnet",
					Labels:    v1alpha1.New().Metadata.Labels,
				},
				Spec: v1alpha1.Spec{
					APIServer: "https://localhost",
					Namespace: "withspec",
				},
				Data: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata":   map[string]interface{}{"name": "config", "namespace": "withspec"},
				},
			},
		},

		{
			name:    "inline",
			baseDir: "./testdata/cases/withenv/",
			expected: manifest.List{{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]interface{}{
					"name":      "config",
					"namespace": "withenv",
				},
			}},
			env: &v1alpha1.Environment{
				APIVersion: v1alpha1.New().APIVersion,
				Kind:       v1alpha1.New().Kind,
				Metadata: v1alpha1.Metadata{
					Name:      "withenv",
					Namespace: "cases/withenv/main.jsonnet",
					Labels:    v1alpha1.New().Metadata.Labels,
				},
				Spec: v1alpha1.Spec{
					APIServer: "https://localhost",
					Namespace: "withenv",
				},
				Data: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata":   map[string]interface{}{"name": "config", "namespace": "withenv"},
				},
			},
		},
		{
			name:    "inline-filename",
			baseDir: "./testdata/cases/withenv/main.jsonnet",
			expected: manifest.List{{
				"apiVersion": "v1",
				"kind":       "ConfigMap",
				"metadata": map[string]interface{}{
					"name":      "config",
					"namespace": "withenv",
				},
			}},
			env: &v1alpha1.Environment{
				APIVersion: v1alpha1.New().APIVersion,
				Kind:       v1alpha1.New().Kind,
				Metadata: v1alpha1.Metadata{
					Name:      "withenv",
					Namespace: "cases/withenv/main.jsonnet",
					Labels:    v1alpha1.New().Metadata.Labels,
				},
				Spec: v1alpha1.Spec{
					APIServer: "https://localhost",
					Namespace: "withenv",
				},
				Data: map[string]interface{}{
					"apiVersion": "v1",
					"kind":       "ConfigMap",
					"metadata":   map[string]interface{}{"name": "config", "namespace": "withenv"},
				},
			},
		},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			l, err := Load(test.baseDir, Opts{})
			require.NoError(t, err)

			assert.Equal(t, test.expected, l.Resources)
			assert.Equal(t, test.env, l.Env)
		})
	}
}

func TestLoadSelectEnvironment(t *testing.T) {
	// No match
	_, err := Load("./testdata/cases/multiple-inline-envs", Opts{Name: "no match"})
	assert.EqualError(t, err, "found no matching environments; run 'tk env list ./testdata/cases/multiple-inline-envs' to view available options")

	// Empty options, match all environments
	_, err = Load("./testdata/cases/multiple-inline-envs", Opts{})
	assert.EqualError(t, err, "found multiple Environments in './testdata/cases/multiple-inline-envs'. Use `--name` to select a single one: \n - project1-env1\n - project1-env2\n - project2-env1")

	// Partial match two environments
	_, err = Load("./testdata/cases/multiple-inline-envs", Opts{Name: "env1"})
	assert.EqualError(t, err, "found multiple Environments in './testdata/cases/multiple-inline-envs'. Use `--name` to select a single one: \n - project1-env1\n - project2-env1")

	// Partial match
	result, err := Load("./testdata/cases/multiple-inline-envs", Opts{Name: "project2"})
	assert.NoError(t, err)
	assert.Equal(t, "project2-env1", result.Env.Metadata.Name)

	// Full match
	result, err = Load("./testdata/cases/multiple-inline-envs", Opts{Name: "project1-env1"})
	assert.NoError(t, err)
	assert.Equal(t, "project1-env1", result.Env.Metadata.Name)
}

func TestLoadSelectEnvironmentFullMatchHasPriority(t *testing.T) {
	// `base` matches both `base` and `base-and-more`
	// However, the full match should win
	result, err := Load("./testdata/cases/inline-name-conflict", Opts{Name: "base"})
	assert.NoError(t, err)
	assert.Equal(t, "base", result.Env.Metadata.Name)
}

func TestLoadFailsWhenBothSpecAndInline(t *testing.T) {
	_, err := Load("./testdata/cases/static-and-inline", Opts{Name: "inline"})
	assert.EqualError(t, err, "found a tanka Environment resource. Check that you aren't using a spec.json and inline environments simultaneously")
}
