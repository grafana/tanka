package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHelp(t *testing.T) {
	output := getCmdOutput(t, "/", "tk", "--help")
	require.Contains(t, output, "Usage")
}
