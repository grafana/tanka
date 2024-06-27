package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplyEnvironment(t *testing.T) {
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
	require.NoError(t, os.WriteFile("environments/default/main.jsonnet", []byte(content), 0600))
	runCmd(t, "tk", "apply", "environments/default", "--auto-approve", "always")
	// Now that the configmap should be there, let's verify it
	runCmd(t, "kubectl", "--namespace", "default", "get", "configmap", "demo")
}
