// Package tools provides LLM-callable tools for the Tanka agent.
package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/jsonschema-go/jsonschema"
	adktool "google.golang.org/adk/tool"
	"google.golang.org/adk/tool/functiontool"
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

// ToADKTools converts a slice of Tools to ADK tool.Tool instances for use with
// the ADK runner and llmagent.
func ToADKTools(ts []Tool) ([]adktool.Tool, error) {
	out := make([]adktool.Tool, 0, len(ts))
	for _, t := range ts {
		at, err := toADKTool(t)
		if err != nil {
			return nil, fmt.Errorf("converting tool %q: %w", t.Name, err)
		}
		out = append(out, at)
	}
	return out, nil
}

func toADKTool(t Tool) (adktool.Tool, error) {
	cfg := functiontool.Config{
		Name:        t.Name,
		Description: t.Description,
	}

	// Convert our JSON Schema to *jsonschema.Schema.
	// jsonschema-go supports direct JSON unmarshaling from standard JSON Schema format.
	if len(t.Schema) > 0 {
		var s jsonschema.Schema
		if err := json.Unmarshal(t.Schema, &s); err == nil {
			cfg.InputSchema = &s
		}
	}

	return functiontool.New(cfg, func(ctx adktool.Context, input map[string]any) (map[string]any, error) {
		inputBytes, err := json.Marshal(input)
		if err != nil {
			return nil, fmt.Errorf("marshaling tool input: %w", err)
		}
		output, err := t.Execute(ctx, json.RawMessage(inputBytes))
		if err != nil {
			return nil, err
		}
		return map[string]any{"output": output}, nil
	})
}
