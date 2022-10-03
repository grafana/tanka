package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiffStat(t *testing.T) {
	cases := []string{
		"empty",
		"added-and-removed",
		"changed-attributes",
		"changed-lots-of-attributes",
	}
	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			content, err := os.ReadFile("testdata/" + c + ".diff")
			require.NoError(t, err)
			expected, err := os.ReadFile("testdata/" + c + ".stat")
			require.NoError(t, err)

			got, err := DiffStat(string(content))
			require.NoError(t, err)

			assert.Equal(t, string(expected), got)
		})
	}
}
