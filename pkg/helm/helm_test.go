package helm

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidateReleaseName(t *testing.T) {
	tests := map[string]struct {
		input       string
		expectError bool
	}{
		"valid-name": {
			input:       "test",
			expectError: false,
		},
		"path-like": {
			input:       "./hello-somewhere",
			expectError: true,
		},
		"argument-like": {
			input:       "--hello-somewhere",
			expectError: true,
		},
		"with-dot": {
			input:       "this.shouldwork",
			expectError: false,
		},
		"overly-long-name": {
			input:       strings.Repeat("a", 54),
			expectError: true,
		},
		"empty": {
			input:       "",
			expectError: true,
		},
	}

	for testName, test := range tests {
		t.Run(testName, func(t *testing.T) {
			err := validateReleaseName(test.input)
			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
