package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestApplyEnvironment(t *testing.T) {
	tmpDir := t.TempDir()
	runCmd(t, tmpDir, "tk", "init")
	runCmd(t, tmpDir, "tk", "env", "set", "environments/default", "--server=https://kubernetes:6443")
	cm := corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "demo",
		},
	}
	content := fmt.Sprintf(`{config: %s}`, marshalToJSON(t, cm))
	require.NoError(t, os.WriteFile(filepath.Join(tmpDir, "environments/default/main.jsonnet"), []byte(content), 0600))
	runCmd(t, tmpDir, "tk", "apply", "environments/default", "--auto-approve", "always")
	// Now that the configmap should be there, let's verify it
	runCmd(t, tmpDir, "kubectl", "--namespace", "default", "get", "configmap", "demo")
}
