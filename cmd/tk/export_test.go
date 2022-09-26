package main

import (
	"errors"
	"testing"

	"github.com/grafana/tanka/pkg/tanka"
	"github.com/stretchr/testify/assert"
)

func TestDetermineMergeStrategy(t *testing.T) {
	cases := []struct {
		name           string
		deprecatedFlag bool
		mergeStrategy  string
		expected       tanka.ExportMergeStrategy
		expectErr      error
	}{
		{
			name:           "default",
			deprecatedFlag: false,
			mergeStrategy:  "",
			expected:       tanka.ExportMergeStrategyNone,
		},
		{
			name:           "deprecated flag set",
			deprecatedFlag: true,
			expected:       tanka.ExportMergeStrategyFailConflicts,
		},
		{
			name:           "both values set",
			deprecatedFlag: true,
			mergeStrategy:  "fail-conflicts",
			expectErr:      errors.New("cannot use --merge and --merge-strategy at the same time"),
		},
		{
			name:          "fail-conflicts",
			mergeStrategy: "fail-on-conflicts",
			expected:      tanka.ExportMergeStrategyFailConflicts,
		},
		{
			name:          "replace-envs",
			mergeStrategy: "replace-envs",
			expected:      tanka.ExportMergeStrategyReplaceEnvs,
		},
		{
			name:          "bad value",
			mergeStrategy: "blabla",
			expectErr:     errors.New("invalid merge strategy: \"blabla\""),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := determineMergeStrategy(tc.deprecatedFlag, tc.mergeStrategy)
			if tc.expectErr != nil {
				assert.EqualError(t, err, tc.expectErr.Error())
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expected, result)
		})
	}
}
