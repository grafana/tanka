package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAgentCreateDevEnvironment(t *testing.T) {
	// Skip if no LLM API key is available
	if os.Getenv("GEMINI_API_KEY") == "" && os.Getenv("GOOGLE_API_KEY") == "" &&
		os.Getenv("ANTHROPIC_API_KEY") == "" && os.Getenv("OPENAI_API_KEY") == "" {
		t.Skip("skipping agent test: set GEMINI_API_KEY, ANTHROPIC_API_KEY, or OPENAI_API_KEY to run")
	}

	tmpDir := t.TempDir()

	// Initialise a bare git repo (tk agent requires being inside one)
	runCmd(t, tmpDir, "git", "init")
	runCmd(t, tmpDir, "git", "config", "user.email", "test@tanka.dev")
	runCmd(t, tmpDir, "git", "config", "user.name", "Tanka Test")

	// Bootstrap a Tanka project (creates environments/default, installs k8s-libsonnet via jb)
	runCmd(t, tmpDir, "tk", "init", "--force")

	// Run the agent in one-shot mode
	runCmd(t, tmpDir, "tk", "agent", "Create a dev environment")

	// 1. environments/dev/ directory must exist
	envDir := filepath.Join(tmpDir, "environments", "dev")
	require.DirExists(t, envDir)

	// 2. spec.json must exist and be valid JSON
	specPath := filepath.Join(envDir, "spec.json")
	require.FileExists(t, specPath)
	specBytes, err := os.ReadFile(specPath)
	require.NoError(t, err)
	var spec map[string]any
	require.NoError(t, json.Unmarshal(specBytes, &spec), "spec.json must be valid JSON")

	// 3. main.jsonnet must exist
	require.FileExists(t, filepath.Join(envDir, "main.jsonnet"))

	// 4. The environment must render without Jsonnet errors
	runCmd(t, tmpDir, "tk", "show", "--dangerous-allow-redirect", "environments/dev")
}
