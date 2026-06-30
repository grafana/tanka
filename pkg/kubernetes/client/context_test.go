package client

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/objx"
)

func TestTryMSISliceMissingKubeconfigSection(t *testing.T) {
	_, err := tryMSISlice(objx.New(map[string]interface{}{}).Get("clusters"), "clusters")
	if err == nil {
		t.Fatal("expected missing kubeconfig section error")
	}

	var missing ErrorMissingKubeconfigSection
	if !errors.As(err, &missing) {
		t.Fatalf("expected ErrorMissingKubeconfigSection, got %T: %v", err, err)
	}
	if !strings.Contains(err.Error(), "$KUBECONFIG") {
		t.Fatalf("expected KUBECONFIG hint, got %q", err.Error())
	}
}
