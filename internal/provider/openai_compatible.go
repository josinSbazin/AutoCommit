package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/josinSbazin/AutoCommit/internal/config"
)

type OpenAICompatibleProvider struct {
	apiKey   string
	endpoint string
	model    string
}

func NewOpenAICompatibleProvider(cfg *config.Config) (*OpenAICompatibleProvider, error) {
	apiKey := os.Getenv("AUTOCOMMIT_API_KEY")
	if cfg.APIKeyEnv != "" {
		apiKey = os.Getenv(cfg.APIKeyEnv)
	}

	if apiKey == "" {
		commonEnvs := []string{
			"GROQ_API_KEY",
			"TOGETHER_API_KEY",
			"OPENROUTER_API_KEY",
			"FIREWORKS_API_KEY",
		}
		for _, env := range commonEnvs {
			if v := os.Getenv(env); v != "" {
				apiKey = v
				break
			}
		}
	}

	endpoint := cfg.Endpoint
	if endpoint == "" {
		if os.Getenv("GROQ_API_KEY") != "" {
			endpoint = "https://api.groq.com/openai/v1/chat/completions"
		} else if os.Getenv("TOGETHER_API_KEY") != "" {
			endpoint = "https://api.together.xyz/v1/chat/completions"
		} else {
			endpoint = "http://localhost:1234/v1/chat/completions"
		}
	}

	if len(endpoint) > 0 && endpoint[len(endpoint)-1] != '/' {
		if !hasSubstring(endpoint, "/chat/completions") {
			endpoint = endpoint + "/chat/completions"
		}
	}

	model := cfg.Model
	if model == "" {
		model = "llama-3.1-70b-versatile"
	}

	return &OpenAICompatibleProvider{
		apiKey:   apiKey,
		endpoint: endpoint,
		model:    model,
	}, nil
}

func (p *OpenAICompatibleProvider) Name() string {
	return "openai-compatible"
}

func (p *OpenAICompatibleProvider) Validate() error {
	return nil
}

func (p *OpenAICompatibleProvider) Generate(ctx context.Context, prompt string) (string, error) {
	reqBody := map[string]any{
		"model": p.model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"max_tokens":  1024,
		"temperature": 0.3,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", p.endpoint, bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if p.apiKey != "" {
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("empty response from API")
	}

	return result.Choices[0].Message.Content, nil
}

func hasSubstring(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
