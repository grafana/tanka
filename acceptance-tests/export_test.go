package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExportEnvironment(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		tmpDir := t.TempDir()
		runCmd(t, tmpDir, "tk", "init")
		runCmd(t, tmpDir, "tk", "env", "set", "environments/default", "--server=https://kubernetes:6443")
		content := `
	{
		config: {
	         apiVersion: "v1",
	         kind: "ConfigMap",
	         metadata : {
	              name: "demo",
	         },
	         data: {},
		},
	}
`
		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "environments/default/main.jsonnet"), []byte(content), 0600))
		runCmd(t, tmpDir, "tk", "export", "export", "environments/default")
		require.FileExists(t, filepath.Join(tmpDir, "export/v1.ConfigMap-demo.yaml"))
	})

	t.Run("only-labeled", func(t *testing.T) {
		tmpDir := t.TempDir()
		runCmd(t, tmpDir, "tk", "init", "--inline")

		// We have two environments stored here but we only want the one with
		// the label "wanted" set to true:
		content := `
	{
		environment(name):: {
			apiVersion: 'tanka.dev/v1alpha1',
			kind: 'Environment',
			metadata: {
			  name: 'environment/%s' % (name),
			  labels: {
				  'wanted': (if name == 'wanted' then 'true' else 'false'),
			  },
			},
			spec: {
			  namespace: 'test-%s' % (name),
			  inline: "true",
			},
			data: {
				config: {
					 apiVersion: "v1",
					 kind: "ConfigMap",
					 metadata : {
						  name: "demo",
					 },
					 data: {},
				 },
			},
		},
		envs: [$.environment('wanted'), $.environment('unwanted')],
	}
`
		require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "environments/default/main.jsonnet"), []byte(content), 0600))
		runCmd(t, tmpDir, "tk", "export", "--recursive", "-l", "wanted=true", "--format", "{{ .metadata.namespace }}/{{.apiVersion}}.{{.kind}}-{{or .metadata.name .metadata.generateName}}", "export", "environments/default")
		assert.FileExists(t, filepath.Join(tmpDir, "export/test-wanted/v1.ConfigMap-demo.yaml"))
		assert.NoFileExists(t, filepath.Join(tmpDir, "export/test-unwanted/v1.ConfigMap-demo-unwanted.yaml"))
	})
}
