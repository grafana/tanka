package agent

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

// NewModel constructs the appropriate model.LLM backend based on the config.
func NewModel(ctx context.Context, cfg *Config) (model.LLM, error) {
	switch cfg.Provider {
	case ProviderGemini:
		return newGeminiModel(ctx, cfg)
	case ProviderAnthropic:
		return NewAnthropicModel(cfg.APIKeyForProvider(), cfg.Model), nil
	case ProviderOpenAI:
		return NewOpenAIModel(cfg.APIKeyForProvider(), cfg.Model), nil
	default:
		return nil, fmt.Errorf("unknown provider %q: must be one of: gemini, anthropic, openai", cfg.Provider)
	}
}
