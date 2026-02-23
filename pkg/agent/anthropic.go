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

const anthropicAPIURL = "https://api.anthropic.com/v1/messages"
const anthropicVersion = "2023-06-01"

// AnthropicProvider calls the Anthropic Messages API directly over HTTP.
type AnthropicProvider struct {
	apiKey string
	model  string
	client *http.Client
}

// NewAnthropicProvider creates a provider that uses the Anthropic API.
func NewAnthropicProvider(apiKey, model string) *AnthropicProvider {
	return &AnthropicProvider{
		apiKey: apiKey,
		model:  model,
		client: &http.Client{},
	}
}

func (p *AnthropicProvider) Name() string { return ProviderAnthropic }

// anthropicContentBlock is a single content block in an Anthropic message.
type anthropicContentBlock struct {
	Type      string          `json:"type"`
	Text      string          `json:"text,omitempty"`
	ID        string          `json:"id,omitempty"`
	Name      string          `json:"name,omitempty"`
	Input     json.RawMessage `json:"input,omitempty"`
	ToolUseID string          `json:"tool_use_id,omitempty"`
	Content   string          `json:"content,omitempty"`
	IsError   bool            `json:"is_error,omitempty"`
}

// anthropicMessage is a message in the Anthropic conversation format.
type anthropicMessage struct {
	Role    string                  `json:"role"`
	Content []anthropicContentBlock `json:"content"`
}

// anthropicTool is a tool definition for the Anthropic API.
type anthropicTool struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"input_schema"`
}

// anthropicRequest is the request body for the Anthropic Messages API.
type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system"`
	Messages  []anthropicMessage `json:"messages"`
	Tools     []anthropicTool    `json:"tools,omitempty"`
}

// anthropicResponse is the response from the Anthropic Messages API.
type anthropicResponse struct {
	ID         string                  `json:"id"`
	Type       string                  `json:"type"`
	Role       string                  `json:"role"`
	Content    []anthropicContentBlock `json:"content"`
	StopReason string                  `json:"stop_reason"`
	Error      *struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error,omitempty"`
}

func (p *AnthropicProvider) Chat(ctx context.Context, systemPrompt string, messages []Message, toolDefs []tools.Tool) (*Message, error) {
	// Convert messages to Anthropic format
	anthropicMsgs := make([]anthropicMessage, 0, len(messages))
	for _, msg := range messages {
		blocks := make([]anthropicContentBlock, 0, len(msg.Content))
		for _, c := range msg.Content {
			switch c.Type {
			case ContentTypeText:
				blocks = append(blocks, anthropicContentBlock{
					Type: "text",
					Text: c.Text,
				})
			case ContentTypeToolUse:
				blocks = append(blocks, anthropicContentBlock{
					Type:  "tool_use",
					ID:    c.ID,
					Name:  c.Name,
					Input: c.Input,
				})
			case ContentTypeToolResult:
				blocks = append(blocks, anthropicContentBlock{
					Type:      "tool_result",
					ToolUseID: c.ToolUseID,
					Content:   c.Text,
					IsError:   c.IsError,
				})
			}
		}
		anthropicMsgs = append(anthropicMsgs, anthropicMessage{
			Role:    string(msg.Role),
			Content: blocks,
		})
	}

	// Convert tool definitions
	anthropicTools := make([]anthropicTool, 0, len(toolDefs))
	for _, t := range toolDefs {
		anthropicTools = append(anthropicTools, anthropicTool{
			Name:        t.Name,
			Description: t.Description,
			InputSchema: t.Schema,
		})
	}

	reqBody := anthropicRequest{
		Model:     p.model,
		MaxTokens: 8192,
		System:    systemPrompt,
		Messages:  anthropicMsgs,
		Tools:     anthropicTools,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshalling request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, anthropicAPIURL, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", anthropicVersion)
	req.Header.Set("content-type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("calling Anthropic API: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var errResp anthropicResponse
		if json.Unmarshal(respBytes, &errResp) == nil && errResp.Error != nil {
			return nil, fmt.Errorf("anthropic API error %d: %s: %s", resp.StatusCode, errResp.Error.Type, errResp.Error.Message)
		}
		return nil, fmt.Errorf("anthropic API returned HTTP %d: %s", resp.StatusCode, string(respBytes))
	}

	var apiResp anthropicResponse
	if err := json.Unmarshal(respBytes, &apiResp); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	// Convert response to internal format
	result := &Message{Role: RoleAssistant}
	for _, block := range apiResp.Content {
		switch block.Type {
		case "text":
			result.Content = append(result.Content, Content{
				Type: ContentTypeText,
				Text: block.Text,
			})
		case "tool_use":
			result.Content = append(result.Content, Content{
				Type:  ContentTypeToolUse,
				ID:    block.ID,
				Name:  block.Name,
				Input: block.Input,
			})
		}
	}

	return result, nil
}
