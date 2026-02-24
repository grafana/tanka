package agent

import (
	"context"
	"fmt"

	"google.golang.org/adk/model"

	"github.com/grafana/tanka/pkg/agent/models"
)

// NewModel constructs the appropriate model.LLM backend based on the config.
func NewModel(ctx context.Context, cfg *Config) (model.LLM, error) {
	apiKey := cfg.APIKeyForProvider()
	modelID := cfg.Model
	switch cfg.Provider {
	case models.ProviderGemini:
		return models.NewGeminiModel(ctx, apiKey, modelID)
	case models.ProviderAnthropic:
		return models.NewAnthropicModel(apiKey, modelID), nil
	case models.ProviderOpenAI:
		return models.NewOpenAIModel(apiKey, modelID), nil
	default:
		return nil, fmt.Errorf("unknown provider %q: must be one of: gemini, anthropic, openai", cfg.Provider)
	}
}
