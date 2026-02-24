// Package tools provides LLM-callable tools for the Tanka agent.
package tools

import (
	"encoding/json"
	"fmt"

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
func bind(input map[string]any, out any) error {
	b, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("marshaling input: %w", err)
	}
	return json.Unmarshal(b, out)
}
