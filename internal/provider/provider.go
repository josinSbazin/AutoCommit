package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/josinSbazin/AutoCommit/internal/config"
)

type Provider interface {
	Generate(ctx context.Context, prompt string) (string, error)
	Name() string
	Validate() error
}

func Get(cfg *config.Config) (Provider, error) {
	switch cfg.Provider {
	case "anthropic", "claude":
		return NewAnthropicProvider(cfg)
	case "openai", "gpt":
		return NewOpenAIProvider(cfg)
	case "ollama":
		return NewOllamaProvider(cfg)
	case "gigachat":
		return NewGigaChatProvider(cfg)
	case "yandexgpt", "yandex":
		return NewYandexGPTProvider(cfg)
	case "openai-compatible":
		return NewOpenAICompatibleProvider(cfg)
	case "":
		return AutoDetect(cfg)
	default:
		return nil, fmt.Errorf("unknown provider: %s", cfg.Provider)
	}
}

func AutoDetect(cfg *config.Config) (Provider, error) {
	checks := []struct {
		envKey   string
		provider string
	}{
		{"ANTHROPIC_API_KEY", "anthropic"},
		{"OPENAI_API_KEY", "openai"},
		{"GIGACHAT_CLIENT_ID", "gigachat"},
		{"YANDEX_API_KEY", "yandexgpt"},
		{"GOOGLE_API_KEY", "google"},
		{"MISTRAL_API_KEY", "mistral"},
		{"GROQ_API_KEY", "openai-compatible"},
	}

	for _, check := range checks {
		if os.Getenv(check.envKey) != "" {
			cfg.Provider = check.provider
			return Get(cfg)
		}
	}

	if isOllamaRunning() {
		cfg.Provider = "ollama"
		return NewOllamaProvider(cfg)
	}

	return nil, fmt.Errorf("no provider configured. Run 'autocommit init' or set API key")
}

func isOllamaRunning() bool {
	return os.Getenv("OLLAMA_HOST") != ""
}
