package client

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/objx"
)

func TestTryMSISlice(t *testing.T) {
	// Test when v is not nil and holds a valid slice
	validVal := objx.New([]map[string]interface{}{{"name": "test"}})
	res, err := tryMSISlice(validVal, "contexts")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(res) != 1 || res[0]["name"] != "test" {
		t.Fatalf("unexpected result: %v", res)
	}

	// Test when v.Data() is nil (triggering makeKubeconfigError)
	nilVal := objx.New(nil)

	// Case 1: KUBECONFIG is set but file doesn't exist
	t.Run("KUBECONFIG set with missing file", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "nonexistent-kubeconfig")
		t.Setenv("KUBECONFIG", tempFile)

		_, err := tryMSISlice(nilVal, "clusters")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		expectedMsg := "The following files in your KUBECONFIG environment variable do not exist"
		if !strings.Contains(err.Error(), expectedMsg) {
			t.Errorf("expected error message to contain '%s', got '%s'", expectedMsg, err.Error())
		}
	})

	// Case 2: KUBECONFIG is unset and default config is missing
	t.Run("KUBECONFIG unset, default config missing", func(t *testing.T) {
		t.Setenv("KUBECONFIG", "")
		// Temporarily change user home to a tmp dir where .kube/config doesn't exist
		tempHome := t.TempDir()

		// Setenv HOME to redirect os.UserHomeDir on Unix systems
		t.Setenv("HOME", tempHome)

		_, err := tryMSISlice(nilVal, "clusters")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		expectedMsg := "KUBECONFIG is unset and the default kubeconfig file at"
		if !strings.Contains(err.Error(), expectedMsg) {
			t.Errorf("expected error to contain '%s', got '%s'", expectedMsg, err.Error())
		}
	})
}
