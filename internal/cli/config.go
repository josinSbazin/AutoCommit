package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/josinSbazin/AutoCommit/internal/config"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View or edit configuration",
	RunE:  runConfig,
}

func init() {
	configCmd.Flags().String("set-provider", "", "Set provider")
	configCmd.Flags().String("set-model", "", "Set model")
	configCmd.Flags().String("set-language", "", "Set language")
	configCmd.Flags().String("set-style", "", "Set style")
}

func runConfig(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	modified := false

	if p, _ := cmd.Flags().GetString("set-provider"); p != "" {
		cfg.Provider = p
		modified = true
	}
	if m, _ := cmd.Flags().GetString("set-model"); m != "" {
		cfg.Model = m
		modified = true
	}
	if l, _ := cmd.Flags().GetString("set-language"); l != "" {
		cfg.Language = l
		modified = true
	}
	if s, _ := cmd.Flags().GetString("set-style"); s != "" {
		cfg.Style = s
		modified = true
	}

	if modified {
		if err := config.Save(cfg, ".autocommit.yml"); err != nil {
			return fmt.Errorf("failed to save config: %w", err)
		}
		fmt.Println("Config updated")
		return nil
	}

	fmt.Println("Current configuration:")
	fmt.Println()

	out, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	fmt.Println(string(out))

	fmt.Println("---")
	fmt.Printf("Sources: %s\n", config.GetLoadedSources())

	return nil
}
