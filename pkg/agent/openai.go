package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/grafana/tanka/pkg/agent/tools"
)

const openAIAPIURL = "https://api.openai.com/v1/chat/completions"

// OpenAIProvider calls the OpenAI Chat Completions API directly over HTTP.
type OpenAIProvider struct {
	apiKey string
	model  string
	client *http.Client
}

// NewOpenAIProvider creates a provider that uses the OpenAI API.
func NewOpenAIProvider(apiKey, model string) *OpenAIProvider {
	return &OpenAIProvider{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{},
	}
}

func (p *OpenAIProvider) Name() string { return ProviderOpenAI }

// openAIMessage is a message in the OpenAI conversation format.
type openAIMessage struct {
	Role       string           `json:"role"`
	Content    string           `json:"content,omitempty"`
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

// openAITool is a tool definition for the OpenAI API.
type openAITool struct {
	Type     string `json:"type"`
	Function struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Parameters  json.RawMessage `json:"parameters"`
	} `json:"function"`
}

// openAIRequest is the request body for the OpenAI Chat Completions API.
type openAIRequest struct {
	Model    string          `json:"model"`
	Messages []openAIMessage `json:"messages"`
	Tools    []openAITool    `json:"tools,omitempty"`
}

// openAIResponse is the response from the OpenAI Chat Completions API.
type openAIResponse struct {
	Choices []struct {
		Message      openAIMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func (p *OpenAIProvider) Chat(ctx context.Context, systemPrompt string, messages []Message, toolDefs []tools.Tool) (*Message, error) {
	// Build OpenAI message list, starting with system message
	openAIMsgs := []openAIMessage{
		{Role: "system", Content: systemPrompt},
	}

	for _, msg := range messages {
		switch msg.Role {
		case RoleUser:
			// User messages: text content and tool results are separate messages
			for _, c := range msg.Content {
				switch c.Type {
				case ContentTypeText:
					openAIMsgs = append(openAIMsgs, openAIMessage{
						Role:    "user",
						Content: c.Text,
					})
				case ContentTypeToolResult:
					openAIMsgs = append(openAIMsgs, openAIMessage{
						Role:       "tool",
						Content:    c.Text,
						ToolCallID: c.ToolUseID,
					})
				}
			}
		case RoleAssistant:
			// Assistant messages: may contain text and/or tool calls
			var textContent string
			var toolCalls []openAIToolCall
			for _, c := range msg.Content {
				switch c.Type {
				case ContentTypeText:
					textContent = c.Text
				case ContentTypeToolUse:
					toolCalls = append(toolCalls, openAIToolCall{
						ID:   c.ID,
						Type: "function",
						Function: struct {
							Name      string `json:"name"`
							Arguments string `json:"arguments"`
						}{
							Name:      c.Name,
							Arguments: string(c.Input),
						},
					})
				}
			}
			openAIMsgs = append(openAIMsgs, openAIMessage{
				Role:      "assistant",
				Content:   textContent,
				ToolCalls: toolCalls,
			})
		}
	}

	// Convert tool definitions
	openAITools := make([]openAITool, 0, len(toolDefs))
	for _, t := range toolDefs {
		var tool openAITool
		tool.Type = "function"
		tool.Function.Name = t.Name
		tool.Function.Description = t.Description
		tool.Function.Parameters = t.Schema
		openAITools = append(openAITools, tool)
	}

	reqBody := openAIRequest{
		Model:    p.model,
		Messages: openAIMsgs,
		Tools:    openAITools,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, openAIAPIURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
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

	choice := apiResp.Choices[0].Message
	result := &Message{Role: RoleAssistant}

	if choice.Content != "" {
		result.Content = append(result.Content, Content{
			Type: ContentTypeText,
			Text: choice.Content,
		})
	}
	for _, tc := range choice.ToolCalls {
		result.Content = append(result.Content, Content{
			Type:  ContentTypeToolUse,
			ID:    tc.ID,
			Name:  tc.Function.Name,
			Input: []byte(tc.Function.Arguments),
		})
	}

	return result, nil
}
