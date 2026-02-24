// Package tools provides LLM-callable tools for the Tanka agent.
package tools

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/google/jsonschema-go/jsonschema"
)

// mustSchema parses a JSON Schema string into *jsonschema.Schema.
// Panics if the JSON is malformed — schemas are hardcoded constants.
func mustSchema(raw string) *jsonschema.Schema {
	var s jsonschema.Schema
	if err := json.Unmarshal([]byte(raw), &s); err != nil {
		panic(fmt.Sprintf("invalid tool schema: %v", err))
	}
	return &s
}

// bind JSON-round-trips a map[string]any into a typed struct, making it easy
// to extract tool parameters with proper type handling (including arrays).
// Fields tagged with `aliases:"a,b"` are resolved: if the primary json key is
// absent from input, the first matching alias value is promoted to it.
func bind(input map[string]any, out any) error {
	t := reflect.TypeOf(out)
	if t != nil && t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
		t = t.Elem()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			primary := field.Tag.Get("json")
			if idx := strings.Index(primary, ","); idx != -1 {
				primary = primary[:idx]
			}
			if primary == "" || primary == "-" {
				continue
			}
			if _, ok := input[primary]; ok {
				continue
			}
			for _, alias := range strings.Split(field.Tag.Get("aliases"), ",") {
				alias = strings.TrimSpace(alias)
				if alias == "" {
					continue
				}
				if v, ok := input[alias]; ok {
					input[primary] = v
					break
				}
			}
		}
	}
	b, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("marshaling input: %w", err)
	}
	return json.Unmarshal(b, out)
}
