package models

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

// AnthropicModel implements model.LLM against the Anthropic Messages API.
type AnthropicModel struct {
	client anthropic.Client
	model  anthropic.Model
}

// NewAnthropicModel creates a model.LLM backed by the Anthropic API.
func NewAnthropicModel(apiKey, modelID string) *AnthropicModel {
	return &AnthropicModel{
		client: anthropic.NewClient(option.WithAPIKey(apiKey)),
		model:  anthropic.Model(modelID),
	}
}

func (m *AnthropicModel) Name() string { return ProviderAnthropic }

func (m *AnthropicModel) GenerateContent(ctx context.Context, req *model.LLMRequest, _ bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		resp, err := m.call(ctx, req)
		yield(resp, err)
	}
}

func (m *AnthropicModel) call(ctx context.Context, req *model.LLMRequest) (*model.LLMResponse, error) {
	// 1. Extract system prompt
	var systemBlocks []anthropic.TextBlockParam
	if req.Config != nil && req.Config.SystemInstruction != nil {
		for _, part := range req.Config.SystemInstruction.Parts {
			if part.Text != "" {
				systemBlocks = append(systemBlocks, anthropic.TextBlockParam{Text: part.Text})
			}
		}
	}

	// 2. Convert conversation history
	var msgs []anthropic.MessageParam
	for _, content := range req.Contents {
		var blocks []anthropic.ContentBlockParamUnion
		for _, part := range content.Parts {
			switch {
			case part.Text != "":
				blocks = append(blocks, anthropic.NewTextBlock(part.Text))
			case part.FunctionCall != nil:
				blocks = append(blocks, anthropic.NewToolUseBlock(part.FunctionCall.ID, part.FunctionCall.Args, part.FunctionCall.Name))
			case part.FunctionResponse != nil:
				respJSON, _ := json.Marshal(part.FunctionResponse.Response)
				blocks = append(blocks, anthropic.NewToolResultBlock(part.FunctionResponse.ID, string(respJSON), false))
			}
		}
		if len(blocks) == 0 {
			continue
		}
		if content.Role == string(genai.RoleModel) {
			msgs = append(msgs, anthropic.NewAssistantMessage(blocks...))
		} else {
			msgs = append(msgs, anthropic.NewUserMessage(blocks...))
		}
	}

	// 3. Convert tool declarations
	var tools []anthropic.ToolUnionParam
	if req.Config != nil {
		for _, t := range req.Config.Tools {
			for _, decl := range t.FunctionDeclarations {
				inputSchema := anthropic.ToolInputSchemaParam{
					Type: "object",
				}
				if decl.Parameters != nil {
					if len(decl.Parameters.Properties) > 0 {
						props := map[string]any{}
						for k, v := range decl.Parameters.Properties {
							props[k] = genaiSchemaToMap(v)
						}
						inputSchema.Properties = props
					}
					inputSchema.Required = decl.Parameters.Required
				}
				tools = append(tools, anthropic.ToolUnionParam{
					OfTool: &anthropic.ToolParam{
						Name:        decl.Name,
						Description: anthropic.String(decl.Description),
						InputSchema: inputSchema,
					},
				})
			}
		}
	}

	// 4. Call the API
	params := anthropic.MessageNewParams{
		Model:     m.model,
		MaxTokens: 8192,
		Messages:  msgs,
		System:    systemBlocks,
		Tools:     tools,
	}

	message, err := m.client.Messages.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("calling anthropic API: %w", err)
	}

	// 5. Convert response to genai.Content
	var parts []*genai.Part
	for _, block := range message.Content {
		switch block.Type {
		case "text":
			if block.Text != "" {
				parts = append(parts, genai.NewPartFromText(block.Text))
			}
		case "tool_use":
			var args map[string]any
			if err := json.Unmarshal(block.Input, &args); err != nil {
				args = map[string]any{"_error": fmt.Sprintf("failed to parse tool arguments: %v", err)}
			}
			p := genai.NewPartFromFunctionCall(block.Name, args)
			p.FunctionCall.ID = block.ID
			parts = append(parts, p)
		}
	}

	isFinal := message.StopReason != "tool_use"
	return &model.LLMResponse{
		Content: &genai.Content{
			Role:  string(genai.RoleModel),
			Parts: parts,
		},
		TurnComplete: isFinal,
	}, nil
}

// genaiSchemaToMap converts a genai.Schema to a JSON Schema-compatible map.
// Gemini uses uppercase type names (e.g. "OBJECT"); JSON Schema uses lowercase.
func genaiSchemaToMap(s *genai.Schema) map[string]any {
	if s == nil {
		return map[string]any{"type": "object"}
	}
	m := map[string]any{}
	if t := strings.ToLower(string(s.Type)); t != "" {
		m["type"] = t
	}
	if s.Description != "" {
		m["description"] = s.Description
	}
	if len(s.Properties) > 0 {
		props := map[string]any{}
		for k, v := range s.Properties {
			props[k] = genaiSchemaToMap(v)
		}
		m["properties"] = props
	}
	if len(s.Required) > 0 {
		m["required"] = s.Required
	}
	if s.Items != nil {
		m["items"] = genaiSchemaToMap(s.Items)
	}
	if len(s.Enum) > 0 {
		m["enum"] = s.Enum
	}
	return m
}
