package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func runCmd(t *testing.T, dir string, cmd string, args ...string) {
	t.Helper()
	c := exec.Command(cmd, args...)
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	require.NoError(t, err)
}

func getCmdOutput(t *testing.T, dir string, cmd string, args ...string) string {
	t.Helper()
	c := exec.Command(cmd, args...)
	c.Dir = dir
	output, err := c.CombinedOutput()
	require.NoError(t, err)
	return string(output)
}

func marshalToJSON(t *testing.T, obj any) string {
	t.Helper()
	output, err := json.Marshal(obj)
	require.NoError(t, err)
	return string(output)
}
