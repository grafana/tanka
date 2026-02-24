package models

import (
	"context"

	"google.golang.org/adk/model"
	adkgemini "google.golang.org/adk/model/gemini"
	"google.golang.org/genai"
)

// NewGeminiModel creates a model.LLM backed by the Gemini API.
func NewGeminiModel(ctx context.Context, apiKey, modelID string) (model.LLM, error) {
	return adkgemini.NewModel(ctx, modelID, &genai.ClientConfig{
		APIKey: apiKey,
	})
}
