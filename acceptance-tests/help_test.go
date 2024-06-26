package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHelp(t *testing.T) {
	output, err := exec.Command("tk", "--help").CombinedOutput()
	require.NoError(t, err)
	require.Contains(t, string(output), "Usage")
}

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

func runCmd(t *testing.T, cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	require.NoError(t, err)
}
