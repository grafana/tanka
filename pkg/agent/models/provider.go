package models

import (
	"context"
	"fmt"

	"google.golang.org/adk/model"
)

// Provider name constants used across config, provider selection, and Name() methods.
const (
	ProviderGemini    = "gemini"
	ProviderAnthropic = "anthropic"
	ProviderOpenAI    = "openai"
)

// ModelConfig holds the minimal configuration needed to construct an LLM.
type ModelConfig struct {
	Provider string
	Model    string
	APIKey   string
}

// NewModel constructs the appropriate model.LLM backend based on the config.
func NewModel(ctx context.Context, cfg *ModelConfig) (model.LLM, error) {
	switch cfg.Provider {
	case ProviderGemini:
		return NewGeminiModel(ctx, cfg.APIKey, cfg.Model)
	case ProviderAnthropic:
		return NewAnthropicModel(cfg.APIKey, cfg.Model), nil
	case ProviderOpenAI:
		return NewOpenAIModel(cfg.APIKey, cfg.Model), nil
	default:
		return nil, fmt.Errorf("unknown provider %q: must be one of: gemini, anthropic, openai", cfg.Provider)
	}
}
