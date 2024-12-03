package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestJsonnetImplementationFlag(t *testing.T) {
	// The jsonnet implementation flag should be present for the following sub-commands
	supportedSubCommands := []string{
		"eval",
		"export",
		"status",
		"apply",
		"diff",
		"delete",
		"show",
		// https://github.com/grafana/tanka/pull/1208
		"env list",
	}
	tmpDir := t.TempDir()
	for _, subcommand := range supportedSubCommands {
		t.Run(subcommand, func(t *testing.T) {
			args := []string{}
			command := strings.Split(subcommand, " ")
			args = append(args, command...)
			args = append(args, "--help")
			helpOutput := getCmdOutput(t, tmpDir, "tk", args...)
			require.Contains(t, helpOutput, "jsonnet-implementation")
		})
	}
}
