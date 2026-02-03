package config

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration options
type Config struct {
	// LLM Provider settings
	Provider  string `yaml:"provider"`
	Model     string `yaml:"model"`
	Endpoint  string `yaml:"endpoint,omitempty"`
	APIKeyEnv string `yaml:"api_key_env,omitempty"`

	// YandexGPT specific
	FolderID string `yaml:"folder_id,omitempty"`

	// Commit style
	Style    string `yaml:"style"`    // conventional, simple, detailed
	Language string `yaml:"language"` // en, ru, etc.

	// Formatting
	MaxSubjectLength int  `yaml:"max_subject_length"`
	MaxBodyLength    int  `yaml:"max_body_length"`
	IncludeBody      bool `yaml:"include_body"`
	IncludeFooter    bool `yaml:"include_footer"`

	// Conventional Commits
	Conventional ConventionalConfig `yaml:"conventional,omitempty"`

	// Context settings
	Context ContextConfig `yaml:"context"`

	// Behavior
	Behavior BehaviorConfig `yaml:"behavior"`

	// Custom instructions for LLM
	Instructions string `yaml:"instructions,omitempty"`
}

// ConventionalConfig for Conventional Commits style
type ConventionalConfig struct {
	RequireScope  bool     `yaml:"require_scope"`
	AllowedTypes  []string `yaml:"allowed_types,omitempty"`
	AllowedScopes []string `yaml:"allowed_scopes,omitempty"`
}

// ContextConfig for context gathering
type ContextConfig struct {
	IncludeHistory   bool `yaml:"include_history"`
	HistoryCount     int  `yaml:"history_count"`
	IncludeBranch    bool `yaml:"include_branch"`
	IncludeDiffStats bool `yaml:"include_diff_stats"`
}

// BehaviorConfig for runtime behavior
type BehaviorConfig struct {
	AutoStage           bool `yaml:"auto_stage"`
	Interactive         bool `yaml:"interactive"`
	ConfirmBeforeCommit bool `yaml:"confirm_before_commit"`
}

// Default returns default configuration
func Default() *Config {
	return &Config{
		Provider:         "",
		Model:            "",
		Style:            "conventional",
		Language:         "en",
		MaxSubjectLength: 72,
		MaxBodyLength:    500,
		IncludeBody:      true,
		IncludeFooter:    false,
		Conventional: ConventionalConfig{
			RequireScope: false,
			AllowedTypes: []string{"feat", "fix", "docs", "style", "refactor", "test", "chore", "perf", "ci", "build"},
		},
		Context: ContextConfig{
			IncludeHistory:   true,
			HistoryCount:     10,
			IncludeBranch:    true,
			IncludeDiffStats: true,
		},
		Behavior: BehaviorConfig{
			AutoStage:           false,
			Interactive:         true,
			ConfirmBeforeCommit: true,
		},
	}
}

var loadedSources []string

func Load() (*Config, error) {
	loadedSources = []string{}
	cfg := Default()

	if globalPath := getGlobalConfigPath(); globalPath != "" {
		if err := loadFile(globalPath, cfg); err == nil {
			loadedSources = append(loadedSources, globalPath)
		}
	}

	// 2. Load project config
	projectPath := ".autocommit.yml"
	if err := loadFile(projectPath, cfg); err == nil {
		loadedSources = append(loadedSources, projectPath)
	}

	// 3. Load from environment
	loadEnv(cfg)

	return cfg, nil
}

// loadFile loads config from a YAML file
func loadFile(path string, cfg *Config) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(data, cfg)
}

// loadEnv loads config from environment variables
func loadEnv(cfg *Config) {
	if v := os.Getenv("AUTOCOMMIT_PROVIDER"); v != "" {
		cfg.Provider = v
	}
	if v := os.Getenv("AUTOCOMMIT_MODEL"); v != "" {
		cfg.Model = v
	}
	if v := os.Getenv("AUTOCOMMIT_LANGUAGE"); v != "" {
		cfg.Language = v
	}
	if v := os.Getenv("AUTOCOMMIT_STYLE"); v != "" {
		cfg.Style = v
	}
}

// Save saves configuration to a file
func Save(cfg *Config, path string) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func GetLoadedSources() string {
	if len(loadedSources) == 0 {
		return "defaults only"
	}
	return strings.Join(loadedSources, ", ")
}

func getGlobalConfigPath() string {
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "autocommit", "config.yml")
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "autocommit", "config.yml")
}

func GetGlobalConfigDir() string {
	if runtime.GOOS == "windows" {
		if appData := os.Getenv("APPDATA"); appData != "" {
			return filepath.Join(appData, "autocommit")
		}
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "autocommit")
}
