package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShow(t *testing.T) {
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	runCmd(t, "tk", "init")
	runCmd(t, "tk", "env", "set", "environments/default", "--server=https://kubernetes:6443")
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

	expectedOutput := `apiVersion: v1
data: {}
kind: ConfigMap
metadata:
  name: demo
  namespace: default
`
	require.NoError(t, os.WriteFile("environments/default/main.jsonnet", []byte(content), 0600))
	output := getCmdOutput(t, "tk", "show", "--dangerous-allow-redirect", "environments/default")
	require.Equal(t, expectedOutput, output)
}
