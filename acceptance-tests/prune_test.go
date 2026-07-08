package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	initialContent = `{
		config: [
			{ apiVersion: "v1", kind: "Namespace", metadata: { name: "prune-ns-a" } },
			{ apiVersion: "v1", kind: "Namespace", metadata: { name: "prune-ns-b" } },
			{ apiVersion: "v1", kind: "ConfigMap", metadata: { name: "demo-a", namespace: "prune-ns-a" }, data: {} },
			{ apiVersion: "v1", kind: "ConfigMap", metadata: { name: "demo-b", namespace: "prune-ns-b" }, data: {} },
			{ apiVersion: "v1", kind: "ConfigMap", metadata: { name: "demo-orphan", namespace: "prune-ns-b" }, data: {} },
		],
	}`
	updatedContent = `{
		config: [
			{ apiVersion: "v1", kind: "Namespace", metadata: { name: "prune-ns-a" } },
			{ apiVersion: "v1", kind: "Namespace", metadata: { name: "prune-ns-b" } },
			{ apiVersion: "v1", kind: "ConfigMap", metadata: { name: "demo-b", namespace: "prune-ns-b" }, data: {} },
		],
	}`
	nsA = "prune-ns-a"
	nsB = "prune-ns-b"
)

// pruneSetup initializes a new tk environment and applies the initial content, then updates the environment
func pruneSetup(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()
	runCmd(t, tmpDir, "tk", "init")
	runCmd(t, tmpDir, "tk", "env", "set", "environments/default",
		"--server=https://kubernetes:6443",
		"--inject-labels",
	)
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "environments/default/main.jsonnet"), []byte(initialContent), 0600))
	runCmd(t, tmpDir, "tk", "apply", "environments/default", "--auto-approve", "always")

	// assert that the initial content was applied
	runCmd(t, tmpDir, "kubectl", "--namespace", nsA, "get", "configmap", "demo-a")
	runCmd(t, tmpDir, "kubectl", "--namespace", nsB, "get", "configmap", "demo-b")
	runCmd(t, tmpDir, "kubectl", "--namespace", nsB, "get", "configmap", "demo-orphan")

	// update the environment
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "environments/default/main.jsonnet"), []byte(updatedContent), 0600))

	return tmpDir
}

func TestPrune(t *testing.T) {
	tmpDir := pruneSetup(t)

	runCmd(t, tmpDir, "tk", "prune", "environments/default", "--auto-approve", "always")

	runCmdExpectError(t, tmpDir, "kubectl", "--namespace", nsA, "get", "configmap", "demo-a")
	runCmdExpectError(t, tmpDir, "kubectl", "--namespace", nsB, "get", "configmap", "demo-orphan")
	runCmd(t, tmpDir, "kubectl", "--namespace", nsB, "get", "configmap", "demo-b")
}

func TestPruneNamespace(t *testing.T) {
	tmpDir := pruneSetup(t)

	runCmd(t, tmpDir, "tk", "prune", "environments/default", "--namespace", nsA, "--auto-approve", "always")

	runCmdExpectError(t, tmpDir, "kubectl", "--namespace", nsA, "get", "configmap", "demo-a")
	runCmd(t, tmpDir, "kubectl", "--namespace", nsB, "get", "configmap", "demo-b")
	runCmd(t, tmpDir, "kubectl", "--namespace", nsB, "get", "configmap", "demo-orphan")
}
