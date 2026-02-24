// Package agent implements the Tanka AI assistant.
package agent

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the agent configuration loaded from file and environment variables.
type Config struct {
	// Provider is the LLM provider: "gemini", "anthropic", or "openai"
	Provider string `yaml:"provider"`
	// Model is the model identifier (e.g. "gemini-2.0-flash", "claude-opus-4-6", "gpt-4o")
	Model string `yaml:"model"`
	// APIKey optionally stores the API key (prefer environment variables instead)
	APIKey string `yaml:"api_key"`
}

// defaults by provider
var providerDefaults = map[string]string{
	ProviderGemini:    "gemini-2.0-flash",
	ProviderAnthropic: "claude-opus-4-6",
	ProviderOpenAI:    "gpt-4o",
}

// LoadConfig loads configuration from ~/.config/tanka/agent.yaml and overrides
// with environment variables. Precedence: env vars > config file > defaults.
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Provider: ProviderGemini,
		Model:    "gemini-2.0-flash",
	}

	// Load from config file (ignore if missing)
	home, err := os.UserHomeDir()
	if err == nil {
		configPath := filepath.Join(home, ".config", "tanka", "agent.yaml")
		data, readErr := os.ReadFile(configPath)
		if readErr != nil && !os.IsNotExist(readErr) {
			return nil, fmt.Errorf("reading config file %s: %w", configPath, readErr)
		}
		if readErr == nil {
			if err := yaml.Unmarshal(data, cfg); err != nil {
				return nil, fmt.Errorf("parsing config file %s: %w", configPath, err)
			}
		}
	}

	// Environment variable overrides
	if v := os.Getenv("TANKA_AGENT_PROVIDER"); v != "" {
		cfg.Provider = v
	}
	if v := os.Getenv("TANKA_AGENT_MODEL"); v != "" {
		cfg.Model = v
	}

	// Set default model for provider if model wasn't explicitly configured
	if cfg.Model == "" {
		if def, ok := providerDefaults[cfg.Provider]; ok {
			cfg.Model = def
		}
	}

	return cfg, nil
}

// Validate checks that the required API key is present for the selected provider.
func (c *Config) Validate() error {
	switch c.Provider {
	case ProviderGemini:
		if c.APIKey == "" && os.Getenv("GEMINI_API_KEY") == "" && os.Getenv("GOOGLE_API_KEY") == "" {
			return fmt.Errorf("no API key found for provider %q: set GEMINI_API_KEY or GOOGLE_API_KEY environment variable", c.Provider)
		}
	case ProviderAnthropic:
		if c.APIKey == "" && os.Getenv("ANTHROPIC_API_KEY") == "" {
			return fmt.Errorf("no API key found for provider %q: set ANTHROPIC_API_KEY environment variable", c.Provider)
		}
	case ProviderOpenAI:
		if c.APIKey == "" && os.Getenv("OPENAI_API_KEY") == "" {
			return fmt.Errorf("no API key found for provider %q: set OPENAI_API_KEY environment variable", c.Provider)
		}
	default:
		return fmt.Errorf("unknown provider %q: must be one of: gemini, anthropic, openai", c.Provider)
	}
	return nil
}

// APIKeyForProvider returns the resolved API key for the configured provider.
func (c *Config) APIKeyForProvider() string {
	if c.APIKey != "" {
		return c.APIKey
	}
	switch c.Provider {
	case ProviderGemini:
		if k := os.Getenv("GEMINI_API_KEY"); k != "" {
			return k
		}
		return os.Getenv("GOOGLE_API_KEY")
	case ProviderAnthropic:
		return os.Getenv("ANTHROPIC_API_KEY")
	case ProviderOpenAI:
		return os.Getenv("OPENAI_API_KEY")
	}
	return ""
}
