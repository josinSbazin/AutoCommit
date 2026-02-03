package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/josinSbazin/AutoCommit/internal/config"
	"github.com/josinSbazin/AutoCommit/internal/ui"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize autocommit",
	RunE:  runInit,
}

func init() {
	initCmd.Flags().BoolP("global", "g", false, "Initialize globally")
	initCmd.Flags().String("provider", "", "Set provider")
}

func runInit(cmd *cobra.Command, args []string) error {
	global, _ := cmd.Flags().GetBool("global")
	providerFlag, _ := cmd.Flags().GetString("provider")

	fmt.Println("AutoCommit Setup")
	fmt.Println()

	var selectedProvider string
	if providerFlag != "" {
		selectedProvider = providerFlag
	} else {
		selectedProvider = ui.SelectProvider()
	}

	if selectedProvider == "" {
		return fmt.Errorf("no provider selected")
	}

	apiKey, extraConfig := ui.GetProviderCredentials(selectedProvider)

	cfg := config.Default()
	cfg.Provider = selectedProvider

	for k, v := range extraConfig {
		switch k {
		case "folder_id":
			cfg.FolderID = v
		case "model":
			cfg.Model = v
		}
	}

	var configPath string
	if global {
		configDir := config.GetGlobalConfigDir()
		if configDir == "" {
			return fmt.Errorf("failed to determine config directory")
		}
		if err := os.MkdirAll(configDir, 0755); err != nil {
			return fmt.Errorf("failed to create config dir: %w", err)
		}
		configPath = filepath.Join(configDir, "config.yml")
	} else {
		configPath = ".autocommit.yml"
	}

	if err := config.Save(cfg, configPath); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println()
	ui.PrintAPIKeyHint(selectedProvider, apiKey)

	fmt.Println()
	fmt.Printf("Config saved to %s\n", configPath)

	if ui.Confirm("Install git hook?") {
		if err := installHook(); err != nil {
			ui.PrintWarning("Failed to install hook: " + err.Error())
		} else {
			fmt.Println("Git hook installed")
		}
	}

	fmt.Println()
	fmt.Println("Ready! Use 'git commit' or 'autocommit' to generate messages.")

	return nil
}

func installHook() error {
	hookPath := ".git/hooks/prepare-commit-msg"

	if _, err := os.Stat(".git"); os.IsNotExist(err) {
		return fmt.Errorf("not a git repository")
	}

	if err := os.MkdirAll(".git/hooks", 0755); err != nil {
		return err
	}

	hookContent := `#!/bin/sh
COMMIT_MSG_FILE=$1
COMMIT_SOURCE=$2

if [ -n "$COMMIT_SOURCE" ]; then
    exit 0
fi

if [ -s "$COMMIT_MSG_FILE" ]; then
    first_line=$(head -n1 "$COMMIT_MSG_FILE")
    if [ -n "$first_line" ] && ! echo "$first_line" | grep -q "^#"; then
        exit 0
    fi
fi

autocommit generate --hook-mode --output "$COMMIT_MSG_FILE"
`

	return os.WriteFile(hookPath, []byte(hookContent), 0755)
}
