package tools

import (
	"encoding/json"
	"testing"
)

func TestBind(t *testing.T) {
	type simple struct {
		Name  string `json:"name"`
		Count int    `json:"count"`
	}

	type withAliases struct {
		GlobPattern string `json:"glob_pattern" aliases:"glob,pattern"`
		Offset      int    `json:"offset"`
	}

	type multiField struct {
		EnvPath string `json:"env_path" aliases:"path,env"`
		Limit   int    `json:"limit"    aliases:"max"`
	}

	type omitempty struct {
		Value string `json:"value,omitempty" aliases:"val"`
	}

	tests := []struct {
		name    string
		input   map[string]any
		out     any
		want    any
		wantErr bool
	}{
		// ── basic binding, no aliases ──────────────────────────────────────
		{
			name:  "primary key present",
			input: map[string]any{"name": "foo", "count": 3},
			out:   &simple{},
			want:  &simple{Name: "foo", Count: 3},
		},
		{
			name:  "missing optional field defaults to zero",
			input: map[string]any{"name": "bar"},
			out:   &simple{},
			want:  &simple{Name: "bar", Count: 0},
		},

		// ── alias resolution ───────────────────────────────────────────────
		{
			name:  "primary key used when present",
			input: map[string]any{"glob_pattern": "**/*.jsonnet"},
			out:   &withAliases{},
			want:  &withAliases{GlobPattern: "**/*.jsonnet"},
		},
		{
			name:  "first alias promoted when primary absent",
			input: map[string]any{"glob": "**/*.jsonnet"},
			out:   &withAliases{},
			want:  &withAliases{GlobPattern: "**/*.jsonnet"},
		},
		{
			name:  "second alias promoted when primary and first alias absent",
			input: map[string]any{"pattern": "**/*.jsonnet"},
			out:   &withAliases{},
			want:  &withAliases{GlobPattern: "**/*.jsonnet"},
		},
		{
			name:  "primary takes precedence over aliases when both present",
			input: map[string]any{"glob_pattern": "primary", "glob": "alias"},
			out:   &withAliases{},
			want:  &withAliases{GlobPattern: "primary"},
		},
		{
			name:  "first alias wins over second alias",
			input: map[string]any{"glob": "first", "pattern": "second"},
			out:   &withAliases{},
			want:  &withAliases{GlobPattern: "first"},
		},

		// ── multiple fields with aliases ───────────────────────────────────
		{
			name:  "multiple fields resolved via aliases",
			input: map[string]any{"path": "environments/dev", "max": 10},
			out:   &multiField{},
			want:  &multiField{EnvPath: "environments/dev", Limit: 10},
		},
		{
			name:  "second alias used for env_path",
			input: map[string]any{"env": "environments/staging"},
			out:   &multiField{},
			want:  &multiField{EnvPath: "environments/staging"},
		},

		// ── omitempty tag suffix stripped correctly ────────────────────────
		{
			name:  "omitempty suffix does not break alias resolution",
			input: map[string]any{"val": "hello"},
			out:   &omitempty{},
			want:  &omitempty{Value: "hello"},
		},

		// ── no aliases tag ─────────────────────────────────────────────────
		{
			name:  "field with no aliases tag binds normally",
			input: map[string]any{"offset": 5},
			out:   &withAliases{},
			want:  &withAliases{Offset: 5},
		},

		// ── none of primary or aliases present ────────────────────────────
		{
			name:  "field absent with no matching alias yields zero value",
			input: map[string]any{},
			out:   &withAliases{},
			want:  &withAliases{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := bind(tc.input, tc.out)
			if (err != nil) != tc.wantErr {
				t.Fatalf("bind() error = %v, wantErr %v", err, tc.wantErr)
			}
			if err != nil {
				return
			}
			// Compare via json round-trip to avoid reflect.DeepEqual issues with
			// unexported fields or interface types.
			gotJSON := mustMarshal(t, tc.out)
			wantJSON := mustMarshal(t, tc.want)
			if gotJSON != wantJSON {
				t.Errorf("bind() result mismatch\n got:  %s\n want: %s", gotJSON, wantJSON)
			}
		})
	}
}

func mustMarshal(t *testing.T, v any) string {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("json.Marshal(%v): %v", v, err)
	}
	return string(b)
}
