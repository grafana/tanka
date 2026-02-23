package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"iter"
	"net/http"

	"google.golang.org/adk/model"
	"google.golang.org/genai"
)

const openAIAPIURL = "https://api.openai.com/v1/chat/completions"

// OpenAIModel implements model.LLM against the OpenAI Chat Completions API.
type OpenAIModel struct {
	apiKey string
	model  string
	client *http.Client
}

// NewOpenAIModel creates a model.LLM backed by the OpenAI API.
func NewOpenAIModel(apiKey, modelID string) *OpenAIModel {
	return &OpenAIModel{apiKey: apiKey, model: modelID, client: &http.Client{}}
}

func (m *OpenAIModel) Name() string { return ProviderOpenAI }

func (m *OpenAIModel) GenerateContent(ctx context.Context, req *model.LLMRequest, _ bool) iter.Seq2[*model.LLMResponse, error] {
	return func(yield func(*model.LLMResponse, error) bool) {
		resp, err := m.call(ctx, req)
		yield(resp, err)
	}
}

// ---- OpenAI wire types ----

type openAIRequest struct {
	Model    string       `json:"model"`
	Messages []openAIMsg  `json:"messages"`
	Tools    []openAITool `json:"tools,omitempty"`
}

type openAIMsg struct {
	Role       string           `json:"role"`
	Content    *string          `json:"content"` // pointer so we can send null
	ToolCallID string           `json:"tool_call_id,omitempty"`
	ToolCalls  []openAIToolCall `json:"tool_calls,omitempty"`
}

type openAIToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

type openAITool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string         `json:"name"`
		Description string         `json:"description"`
		Parameters  map[string]any `json:"parameters"`
	} `json:"function"`
}

type openAIResponse struct {
	Choices []struct {
		Message      openAIMsg `json:"message"`
		FinishReason string    `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func (m *OpenAIModel) call(ctx context.Context, req *model.LLMRequest) (*model.LLMResponse, error) {
	// 1. Build messages: start with system instruction
	var msgs []openAIMsg
	if req.Config != nil && req.Config.SystemInstruction != nil {
		var system string
		for _, part := range req.Config.SystemInstruction.Parts {
			system += part.Text
		}
		if system != "" {
			msgs = append(msgs, openAIMsg{Role: "system", Content: &system})
		}
	}

	// 2. Convert conversation history
	for _, content := range req.Contents {
		switch content.Role {
		case string(genai.RoleModel):
			var text string
			var toolCalls []openAIToolCall
			for _, part := range content.Parts {
				if part.Text != "" {
					text = part.Text
				}
				if part.FunctionCall != nil {
					argsJSON, _ := json.Marshal(part.FunctionCall.Args)
					tc := openAIToolCall{
						ID:   part.FunctionCall.ID,
						Type: "function",
					}
					tc.Function.Name = part.FunctionCall.Name
					tc.Function.Arguments = string(argsJSON)
					toolCalls = append(toolCalls, tc)
				}
			}
			var contentPtr *string
			if len(toolCalls) == 0 {
				contentPtr = &text
			}
			msgs = append(msgs, openAIMsg{
				Role:      "assistant",
				Content:   contentPtr,
				ToolCalls: toolCalls,
			})

		default: // "user"
			for _, part := range content.Parts {
				if part.Text != "" {
					msgs = append(msgs, openAIMsg{Role: "user", Content: &part.Text})
				}
				if part.FunctionResponse != nil {
					respJSON, _ := json.Marshal(part.FunctionResponse.Response)
					s := string(respJSON)
					msgs = append(msgs, openAIMsg{
						Role:       "tool",
						Content:    &s,
						ToolCallID: part.FunctionResponse.ID,
					})
				}
			}
		}
	}

	// 3. Convert tool declarations
	var tools []openAITool
	if req.Config != nil {
		for _, t := range req.Config.Tools {
			for _, decl := range t.FunctionDeclarations {
				var ot openAITool
				ot.Type = "function"
				ot.Function.Name = decl.Name
				ot.Function.Description = decl.Description
				ot.Function.Parameters = genaiSchemaToMap(decl.Parameters)
				tools = append(tools, ot)
			}
		}
	}

	// 4. Build and send request
	body := openAIRequest{
		Model:    m.model,
		Messages: msgs,
		Tools:    tools,
	}
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIAPIURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+m.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := m.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("calling OpenAI API: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp openAIResponse
		if json.Unmarshal(respBytes, &errResp) == nil && errResp.Error != nil {
			return nil, fmt.Errorf("OpenAI API error %d: %s: %s", resp.StatusCode, errResp.Error.Type, errResp.Error.Message)
		}
		return nil, fmt.Errorf("OpenAI API returned HTTP %d: %s", resp.StatusCode, string(respBytes))
	}

	var apiResp openAIResponse
	if err := json.Unmarshal(respBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}
	if len(apiResp.Choices) == 0 {
		return nil, fmt.Errorf("OpenAI returned no choices")
	}

	// 5. Convert response to genai.Content
	choice := apiResp.Choices[0].Message
	var parts []*genai.Part

	if choice.Content != nil && *choice.Content != "" {
		parts = append(parts, genai.NewPartFromText(*choice.Content))
	}
	for _, tc := range choice.ToolCalls {
		var args map[string]any
		_ = json.Unmarshal([]byte(tc.Function.Arguments), &args)
		p := genai.NewPartFromFunctionCall(tc.Function.Name, args)
		p.FunctionCall.ID = tc.ID
		parts = append(parts, p)
	}

	isFinal := apiResp.Choices[0].FinishReason != "tool_calls"
	return &model.LLMResponse{
		Content: &genai.Content{
			Role:  string(genai.RoleModel),
			Parts: parts,
		},
		TurnComplete: isFinal,
	}, nil
}
