package main

import (
	"errors"
	"testing"

	"github.com/grafana/tanka/pkg/tanka"
	"github.com/stretchr/testify/assert"
)

func TestValidateAutoApprove(t *testing.T) {
	for _, tc := range []struct {
		name                  string
		autoApproveDeprecated bool
		autoApproveString     string
		expected              tanka.AutoApproveSetting
		expectErr             error
	}{
		{
			name:     "default",
			expected: tanka.AutoApproveNever,
		},
		{
			name:                  "deprecated bool set",
			autoApproveDeprecated: true,
			expected:              tanka.AutoApproveAlways,
		},
		{
			name:                  "both values set",
			autoApproveDeprecated: true,
			autoApproveString:     "never",
			expectErr:             errors.New("--dangerous-auto-approve and --auto-approve are mutually exclusive"),
		},
		{
			name:              "always",
			autoApproveString: "always",
			expected:          tanka.AutoApproveAlways,
		},
		{
			name:              "never",
			autoApproveString: "never",
			expected:          tanka.AutoApproveNever,
		},
		{
			name:              "if-no-changes",
			autoApproveString: "if-no-changes",
			expected:          tanka.AutoApproveNoChanges,
		},
		{
			name:              "bad value",
			autoApproveString: "blabla",
			expectErr:         errors.New("invalid value for --auto-approve: blabla"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			result, err := validateAutoApprove(tc.autoApproveDeprecated, tc.autoApproveString)
			if tc.expectErr != nil {
				assert.EqualError(t, err, tc.expectErr.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}
