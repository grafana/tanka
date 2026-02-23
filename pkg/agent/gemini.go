package agent

import (
	"context"

	"google.golang.org/adk/model"
	adkgemini "google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

func newGeminiModel(ctx context.Context, cfg *Config) (model.LLM, error) {
	return adkgemini.NewModel(ctx, cfg.Model, &genai.ClientConfig{
		APIKey: cfg.APIKeyForProvider(),
	})
}
