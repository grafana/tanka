package main

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHelp(t *testing.T) {
	output, err := exec.Command("tk", "--help").CombinedOutput()
	require.NoError(t, err)
	require.Contains(t, string(output), "Usage")
}
