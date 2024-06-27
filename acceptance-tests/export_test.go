package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExportEnvironment(t *testing.T) {
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
}
