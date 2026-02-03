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

type YandexGPTProvider struct {
	apiKey   string
	iamToken string
	folderID string
	model    string
}

func NewYandexGPTProvider(cfg *config.Config) (*YandexGPTProvider, error) {
	apiKey := os.Getenv("YANDEX_API_KEY")
	iamToken := os.Getenv("YANDEX_IAM_TOKEN")
	folderID := os.Getenv("YANDEX_FOLDER_ID")

	if folderID == "" && cfg.FolderID != "" {
		folderID = cfg.FolderID
	}

	if apiKey == "" && iamToken == "" {
		return nil, fmt.Errorf("YANDEX_API_KEY or YANDEX_IAM_TOKEN environment variable not set")
	}

	if folderID == "" {
		return nil, fmt.Errorf("YANDEX_FOLDER_ID environment variable not set")
	}

	model := cfg.Model
	if model == "" {
		model = "yandexgpt"
	}

	return &YandexGPTProvider{
		apiKey:   apiKey,
		iamToken: iamToken,
		folderID: folderID,
		model:    model,
	}, nil
}

func (p *YandexGPTProvider) Name() string {
	return "yandexgpt"
}

func (p *YandexGPTProvider) Validate() error {
	if p.apiKey == "" && p.iamToken == "" {
		return fmt.Errorf("no API key or IAM token configured")
	}
	if p.folderID == "" {
		return fmt.Errorf("folder ID not configured")
	}
	return nil
}

func (p *YandexGPTProvider) Generate(ctx context.Context, prompt string) (string, error) {
	modelURI := fmt.Sprintf("gpt://%s/%s/latest", p.folderID, p.model)

	reqBody := map[string]any{
		"modelUri": modelURI,
		"completionOptions": map[string]any{
			"stream":      false,
			"temperature": 0.3,
			"maxTokens":   "1024",
		},
		"messages": []map[string]string{
			{"role": "user", "text": prompt},
		},
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://llm.api.cloud.yandex.net/foundationModels/v1/completion",
		bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-folder-id", p.folderID)

	if p.apiKey != "" {
		req.Header.Set("Authorization", "Api-Key "+p.apiKey)
	} else {
		req.Header.Set("Authorization", "Bearer "+p.iamToken)
	}

	resp, err := http.DefaultClient.Do(req)
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
		Result struct {
			Alternatives []struct {
				Message struct {
					Role string `json:"role"`
					Text string `json:"text"`
				} `json:"message"`
			} `json:"alternatives"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Result.Alternatives) == 0 {
		return "", fmt.Errorf("empty response from API")
	}

	return result.Result.Alternatives[0].Message.Text, nil
}
