package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"iter"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/shared"
	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

// OpenAIModel implements model.LLM against the OpenAI Chat Completions API.
type OpenAIModel struct {
	client openai.Client
	model  shared.ChatModel
}

// NewOpenAIModel creates a model.LLM backed by the OpenAI API.
func NewOpenAIModel(apiKey, modelID string) *OpenAIModel {
	return &OpenAIModel{
		client: openai.NewClient(option.WithAPIKey(apiKey)),
		model:  shared.ChatModel(modelID),
	}
}

func (m *OpenAIModel) Name() string { return ProviderOpenAI }

func (m *OpenAIModel) GenerateContent(ctx context.Context, req *model.LLMRequest, _ bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		resp, err := m.call(ctx, req)
		yield(resp, err)
	}
}

func (m *OpenAIModel) call(ctx context.Context, req *model.LLMRequest) (*model.LLMResponse, error) {
	// 1. Build messages: start with system instruction
	var msgs []openai.ChatCompletionMessageParamUnion
	if req.Config != nil && req.Config.SystemInstruction != nil {
		var system string
		for _, part := range req.Config.SystemInstruction.Parts {
			system += part.Text
		}
		if system != "" {
			msgs = append(msgs, openai.SystemMessage(system))
		}
	}

	// 2. Convert conversation history
	for _, content := range req.Contents {
		switch content.Role {
		case string(genai.RoleModel):
			var text string
			var toolCalls []openai.ChatCompletionMessageToolCallParam
			for _, part := range content.Parts {
				if part.Text != "" {
					text = part.Text
				}
				if part.FunctionCall != nil {
					argsJSON, _ := json.Marshal(part.FunctionCall.Args)
					toolCalls = append(toolCalls, openai.ChatCompletionMessageToolCallParam{
						ID: part.FunctionCall.ID,
						Function: openai.ChatCompletionMessageToolCallFunctionParam{
							Name:      part.FunctionCall.Name,
							Arguments: string(argsJSON),
						},
					})
				}
			}
			p := openai.ChatCompletionAssistantMessageParam{ToolCalls: toolCalls}
			if len(toolCalls) == 0 {
				p.Content.OfString = openai.String(text)
			}
			msgs = append(msgs, openai.ChatCompletionMessageParamUnion{OfAssistant: &p})

		default: // "user"
			for _, part := range content.Parts {
				if part.Text != "" {
					msgs = append(msgs, openai.UserMessage(part.Text))
				}
				if part.FunctionResponse != nil {
					respJSON, _ := json.Marshal(part.FunctionResponse.Response)
					msgs = append(msgs, openai.ToolMessage(string(respJSON), part.FunctionResponse.ID))
				}
			}
		}
	}

	// 3. Convert tool declarations
	var tools []openai.ChatCompletionToolParam
	if req.Config != nil {
		for _, t := range req.Config.Tools {
			for _, decl := range t.FunctionDeclarations {
				tools = append(tools, openai.ChatCompletionToolParam{
					Function: shared.FunctionDefinitionParam{
						Name:        decl.Name,
						Description: openai.String(decl.Description),
						Parameters:  shared.FunctionParameters(genaiSchemaToMap(decl.Parameters)),
					},
				})
			}
		}
	}

	// 4. Call the API
	resp, err := m.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model:    m.model,
		Messages: msgs,
		Tools:    tools,
	})
	if err != nil {
		return nil, fmt.Errorf("calling OpenAI API: %w", err)
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("OpenAI returned no choices")
	}

	// 5. Convert response to genai.Content
	choice := resp.Choices[0]
	var parts []*genai.Part

	if choice.Message.Content != "" {
		parts = append(parts, genai.NewPartFromText(choice.Message.Content))
	}
	for _, tc := range choice.Message.ToolCalls {
		var args map[string]any
		_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)
		p := genai.NewPartFromFunctionCall(tc.Function.Name, args)
		p.FunctionCall.ID = tc.ID
		parts = append(parts, p)
	}

	isFinal := choice.FinishReason != "tool_calls"
	return &model.LLMResponse{
		Content: &genai.Content{
			Role:  string(genai.RoleModel),
			Parts: parts,
		},
		TurnComplete: isFinal,
	}, nil
}
