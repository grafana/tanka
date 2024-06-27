package main

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func runCmd(t *testing.T, cmd string, args ...string) {
	c := exec.Command(cmd, args...)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	err := c.Run()
	require.NoError(t, err)
}

func getCmdOutput(t *testing.T, cmd string, args ...string) string {
	c := exec.Command(cmd, args...)
	output, err := c.CombinedOutput()
	require.NoError(t, err)
	return string(output)
}
