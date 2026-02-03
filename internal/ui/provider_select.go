package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Provider info for selection
type providerInfo struct {
	name        string
	displayName string
	description string
	envVars     []string
}

var providers = []providerInfo{
	{
		name:        "anthropic",
		displayName: "Anthropic (Claude)",
		description: "Claude 3.5/4 - excellent for code understanding",
		envVars:     []string{"ANTHROPIC_API_KEY"},
	},
	{
		name:        "openai",
		displayName: "OpenAI (GPT-4)",
		description: "GPT-4o - fast and capable",
		envVars:     []string{"OPENAI_API_KEY"},
	},
	{
		name:        "gigachat",
		displayName: "GigaChat (Сбер)",
		description: "Российская модель от Сбера",
		envVars:     []string{"GIGACHAT_CLIENT_ID", "GIGACHAT_CLIENT_SECRET"},
	},
	{
		name:        "yandexgpt",
		displayName: "YandexGPT",
		description: "Российская модель от Яндекса",
		envVars:     []string{"YANDEX_API_KEY", "YANDEX_FOLDER_ID"},
	},
	{
		name:        "ollama",
		displayName: "Ollama (Local)",
		description: "Run models locally - free, private",
		envVars:     []string{},
	},
	{
		name:        "openai-compatible",
		displayName: "OpenAI-Compatible",
		description: "Groq, Together AI, LM Studio, etc.",
		envVars:     []string{"AUTOCOMMIT_API_KEY"},
	},
}

// SelectProvider shows provider selection menu
func SelectProvider() string {
	fmt.Println("Select LLM provider:")
	fmt.Println()

	for i, p := range providers {
		fmt.Printf("  %s[%d]%s %s\n", colorBold, i+1, colorReset, p.displayName)
		fmt.Printf("      %s%s%s\n", colorGray, p.description, colorReset)
	}

	fmt.Println()
	fmt.Print("Enter number (1-6): ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Parse selection
	var idx int
	fmt.Sscanf(input, "%d", &idx)

	if idx < 1 || idx > len(providers) {
		return ""
	}

	return providers[idx-1].name
}

// GetProviderCredentials prompts for API credentials
func GetProviderCredentials(providerName string) (string, map[string]string) {
	extra := make(map[string]string)

	switch providerName {
	case "anthropic":
		fmt.Println()
		fmt.Println("Get your API key from: https://console.anthropic.com/")
		fmt.Print("Enter ANTHROPIC_API_KEY: ")
		return readLine(), extra

	case "openai":
		fmt.Println()
		fmt.Println("Get your API key from: https://platform.openai.com/api-keys")
		fmt.Print("Enter OPENAI_API_KEY: ")
		return readLine(), extra

	case "gigachat":
		fmt.Println()
		fmt.Println("Get credentials from: https://developers.sber.ru/")
		fmt.Print("Enter GIGACHAT_CLIENT_ID: ")
		clientID := readLine()
		fmt.Print("Enter GIGACHAT_CLIENT_SECRET: ")
		clientSecret := readLine()
		extra["client_id"] = clientID
		extra["client_secret"] = clientSecret
		return clientID, extra

	case "yandexgpt":
		fmt.Println()
		fmt.Println("Get credentials from: https://console.cloud.yandex.ru/")
		fmt.Print("Enter YANDEX_API_KEY: ")
		apiKey := readLine()
		fmt.Print("Enter YANDEX_FOLDER_ID: ")
		folderID := readLine()
		extra["folder_id"] = folderID
		return apiKey, extra

	case "ollama":
		fmt.Println()
		fmt.Println("Make sure Ollama is running: https://ollama.ai/")
		fmt.Println("Default model: llama3.1")
		fmt.Print("Enter model name (or press Enter for default): ")
		model := readLine()
		if model != "" {
			extra["model"] = model
		}
		return "", extra

	case "openai-compatible":
		fmt.Println()
		fmt.Println("Enter endpoint URL (e.g., https://api.groq.com/openai/v1):")
		fmt.Print("Endpoint: ")
		endpoint := readLine()
		extra["endpoint"] = endpoint
		fmt.Print("API Key (if required): ")
		apiKey := readLine()
		fmt.Print("Model name: ")
		model := readLine()
		extra["model"] = model
		return apiKey, extra
	}

	return "", extra
}

// PrintAPIKeyHint shows how to set API key
func PrintAPIKeyHint(providerName, apiKey string) {
	if apiKey == "" {
		return
	}

	fmt.Println("Add this to your shell profile (.bashrc, .zshrc, etc.):")
	fmt.Println()

	switch providerName {
	case "anthropic":
		fmt.Printf("  export ANTHROPIC_API_KEY='%s'\n", maskKey(apiKey))
	case "openai":
		fmt.Printf("  export OPENAI_API_KEY='%s'\n", maskKey(apiKey))
	case "gigachat":
		fmt.Println("  export GIGACHAT_CLIENT_ID='...'")
		fmt.Println("  export GIGACHAT_CLIENT_SECRET='...'")
	case "yandexgpt":
		fmt.Printf("  export YANDEX_API_KEY='%s'\n", maskKey(apiKey))
		fmt.Println("  export YANDEX_FOLDER_ID='...'")
	}
}

func readLine() string {
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	return strings.TrimSpace(line)
}

func maskKey(key string) string {
	if len(key) <= 8 {
		return "***"
	}
	return key[:4] + "..." + key[len(key)-4:]
}
