package provider

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/josinSbazin/AutoCommit/internal/config"
)

type GigaChatProvider struct {
	clientID     string
	clientSecret string
	model        string

	mu          sync.Mutex
	accessToken string
	expiresAt   time.Time
	client      *http.Client
}

func NewGigaChatProvider(cfg *config.Config) (*GigaChatProvider, error) {
	clientID := os.Getenv("GIGACHAT_CLIENT_ID")
	clientSecret := os.Getenv("GIGACHAT_CLIENT_SECRET")

	if clientID == "" {
		creds := os.Getenv("GIGACHAT_CREDENTIALS")
		if creds != "" {
			decoded, err := base64.StdEncoding.DecodeString(creds)
			if err == nil {
				parts := strings.SplitN(string(decoded), ":", 2)
				if len(parts) == 2 {
					clientID = parts[0]
					clientSecret = parts[1]
				}
			}
		}
	}

	if clientID == "" || clientSecret == "" {
		return nil, fmt.Errorf("GIGACHAT_CLIENT_ID and GIGACHAT_CLIENT_SECRET environment variables not set")
	}

	model := cfg.Model
	if model == "" {
		model = "GigaChat"
	}

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	return &GigaChatProvider{
		clientID:     clientID,
		clientSecret: clientSecret,
		model:        model,
		client:       &http.Client{Transport: transport, Timeout: 60 * time.Second},
	}, nil
}

func (p *GigaChatProvider) Name() string {
	return "gigachat"
}

func (p *GigaChatProvider) Validate() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return p.authorize(ctx)
}

func (p *GigaChatProvider) authorize(ctx context.Context) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.accessToken != "" && time.Now().Before(p.expiresAt.Add(-time.Minute)) {
		return nil
	}

	credentials := base64.StdEncoding.EncodeToString(
		[]byte(p.clientID + ":" + p.clientSecret))

	data := url.Values{}
	data.Set("scope", "GIGACHAT_API_PERS")

	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://ngw.devices.sberbank.ru:9443/api/v2/oauth",
		strings.NewReader(data.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Authorization", "Basic "+credentials)
	req.Header.Set("RqUID", fmt.Sprintf("%d", time.Now().UnixNano()))

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("auth request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read auth response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("auth failed (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresAt   int64  `json:"expires_at"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("failed to parse auth response: %w", err)
	}

	p.accessToken = result.AccessToken
	p.expiresAt = time.UnixMilli(result.ExpiresAt)

	return nil
}

func (p *GigaChatProvider) Generate(ctx context.Context, prompt string) (string, error) {
	if err := p.authorize(ctx); err != nil {
		return "", fmt.Errorf("authorization failed: %w", err)
	}

	reqBody := map[string]any{
		"model": p.model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
		"temperature": 0.3,
		"max_tokens":  1024,
	}

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST",
		"https://gigachat.devices.sberbank.ru/api/v1/chat/completions",
		bytes.NewReader(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.accessToken)

	resp, err := p.client.Do(req)
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
