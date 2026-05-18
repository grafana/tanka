package process

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKindsFor(t *testing.T) {
	cases := []struct {
		name      string
		exprs     []string
		wantKinds []string
		wantOK    bool
	}{
		{
			name:      "no filters",
			exprs:     nil,
			wantKinds: nil,
			wantOK:    false,
		},
		{
			name:      "empty filters",
			exprs:     []string{},
			wantKinds: nil,
			wantOK:    false,
		},
		{
			name:      "single kind with wildcard name",
			exprs:     []string{"statefulset/.*"},
			wantKinds: []string{"statefulset"},
			wantOK:    true,
		},
		{
			name:      "single kind with literal name",
			exprs:     []string{"statefulset/live-store"},
			wantKinds: []string{"statefulset"},
			wantOK:    true,
		},
		{
			name:      "single kind with name regex",
			exprs:     []string{"statefulset/live-store.*"},
			wantKinds: []string{"statefulset"},
			wantOK:    true,
		},
		{
			name:      "multiple kinds",
			exprs:     []string{"statefulset/.*", "deployment/.*"},
			wantKinds: []string{"statefulset", "deployment"},
			wantOK:    true,
		},
		{
			name:      "only negative matcher",
			exprs:     []string{"!statefulset/.*"},
			wantKinds: nil,
			wantOK:    false,
		},
		{
			name:      "positive and negative matcher",
			exprs:     []string{"statefulset/.*", "!deployment/.*"},
			wantKinds: []string{"statefulset"},
			wantOK:    true,
		},
		{
			name:      "kind with regex metachar — falls back to all",
			exprs:     []string{"(statefulset|deployment)/.*"},
			wantKinds: nil,
			wantOK:    false,
		},
		{
			name:      "wildcard kind — falls back to all",
			exprs:     []string{".*/.*"},
			wantKinds: nil,
			wantOK:    false,
		},
		{
			name:      "kind without slash",
			exprs:     []string{"statefulset"},
			wantKinds: []string{"statefulset"},
			wantOK:    true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			matchers, err := StrExps(c.exprs...)
			require.NoError(t, err)

			kinds, ok := matchers.KindsFor()
			assert.Equal(t, c.wantOK, ok)
			assert.Equal(t, c.wantKinds, kinds)
		})
	}
}

func TestKindFromFilterPattern(t *testing.T) {
	cases := []struct {
		name     string
		pattern  string
		wantKind string
		wantOK   bool
	}{
		{
			name:     "statefulset with wildcard",
			pattern:  `(?i)^statefulset/.*$`,
			wantKind: "statefulset",
			wantOK:   true,
		},
		{
			name:     "statefulset with name regex",
			pattern:  `(?i)^statefulset/live-store.*$`,
			wantKind: "statefulset",
			wantOK:   true,
		},
		{
			name:     "kind with no slash",
			pattern:  `(?i)^statefulset$`,
			wantKind: "statefulset",
			wantOK:   true,
		},
		{
			name:     "kind alternation falls back",
			pattern:  `(?i)^(statefulset|deployment)/.*$`,
			wantKind: "",
			wantOK:   false,
		},
		{
			name:     "wildcard kind falls back",
			pattern:  `(?i)^.*/.*$`,
			wantKind: "",
			wantOK:   false,
		},
		{
			name:     "mixed case kind",
			pattern:  `(?i)^StatefulSet/.*$`,
			wantKind: "StatefulSet",
			wantOK:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			kind, ok := kindFromFilterPattern(c.pattern)
			assert.Equal(t, c.wantOK, ok)
			assert.Equal(t, c.wantKind, kind)
		})
	}
}
