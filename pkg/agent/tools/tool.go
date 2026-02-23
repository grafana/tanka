// Package tools provides LLM-callable tools for the Tanka agent.
package tools

import (
	"context"
	"encoding/json"
)

// Tool represents a callable function that the LLM can invoke.
type Tool struct {
	// Name is the identifier the LLM uses to invoke this tool.
	Name string
	// Description explains what the tool does (shown to the LLM).
	Description string
	// Schema is the JSON Schema (object type) describing the input parameters.
	Schema json.RawMessage
	// Execute runs the tool with the given JSON input and returns a string result.
	Execute func(ctx context.Context, input json.RawMessage) (string, error)
}
