package agent

import (
	"context"
	"fmt"

	"github.com/grafana/tanka/pkg/agent/tools"
)

// Provider name constants used across config, provider selection, and Name() methods.
const (
	ProviderGemini    = "gemini"
	ProviderAnthropic = "anthropic"
	ProviderOpenAI    = "openai"
)

// Role identifies the author of a conversation message.
type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

// ContentType identifies the kind of content within a message.
type ContentType string

const (
	// ContentTypeText is a plain text block.
	ContentTypeText ContentType = "text"
	// ContentTypeToolUse is a tool invocation by the assistant.
	ContentTypeToolUse ContentType = "tool_use"
	// ContentTypeToolResult is the result of a tool call, sent by the user.
	ContentTypeToolResult ContentType = "tool_result"
)

// Content is a single block of content within a Message.
type Content struct {
	Type ContentType

	// Text holds the content for ContentTypeText and ContentTypeToolResult.
	Text string

	// ID is the unique identifier for a ContentTypeToolUse block.
	ID string
	// Name is the tool name for a ContentTypeToolUse block.
	Name string
	// Input is the JSON-encoded input for a ContentTypeToolUse block.
	Input []byte

	// ToolUseID references the ContentTypeToolUse ID for a ContentTypeToolResult block.
	ToolUseID string
	// IsError marks a ContentTypeToolResult as an error response.
	IsError bool
}

// Message is a single turn in the conversation.
type Message struct {
	Role    Role
	Content []Content
}

// TextMessage creates a simple user or assistant text message.
func TextMessage(role Role, text string) Message {
	return Message{
		Role:    role,
		Content: []Content{{Type: ContentTypeText, Text: text}},
	}
}

// Provider is the interface that LLM backends must implement.
type Provider interface {
	// Name returns a human-readable identifier for the provider (e.g. "anthropic").
	Name() string
	// Chat sends the conversation to the LLM and returns the assistant's response.
	// The response may contain text, tool calls, or both.
	Chat(ctx context.Context, systemPrompt string, messages []Message, tools []tools.Tool) (*Message, error)
}

// NewProvider constructs the appropriate Provider based on the config.
func NewProvider(cfg *Config) (Provider, error) {
	switch cfg.Provider {
	case ProviderAnthropic:
		return NewAnthropicProvider(cfg.APIKeyForProvider(), cfg.Model), nil
	case ProviderGemini:
		return NewGeminiProvider(cfg.APIKeyForProvider(), cfg.Model), nil
	case ProviderOpenAI:
		return NewOpenAIProvider(cfg.APIKeyForProvider(), cfg.Model), nil
	default:
		return nil, fmt.Errorf("unsupported provider %q", cfg.Provider)
	}
}
