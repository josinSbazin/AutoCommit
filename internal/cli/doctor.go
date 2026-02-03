package cli

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
	"github.com/josinSbazin/AutoCommit/internal/config"
	"github.com/josinSbazin/AutoCommit/internal/git"
	"github.com/josinSbazin/AutoCommit/internal/provider"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose issues",
	RunE:  runDoctor,
}

func runDoctor(cmd *cobra.Command, args []string) error {
	fmt.Println("AutoCommit Doctor")
	fmt.Println()

	issues := 0

	fmt.Print("Git installation... ")
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println("Not found")
		issues++
	} else {
		fmt.Println("OK")
	}

	fmt.Print("Git repository... ")
	if git.IsRepo() {
		fmt.Println("OK")
	} else {
		fmt.Println("Not in a git repository")
	}

	fmt.Print("Configuration... ")
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		issues++
	} else {
		fmt.Println("OK")
	}

	if cfg != nil {
		fmt.Printf("Provider (%s)... ", cfg.Provider)
		if cfg.Provider == "" {
			fmt.Println("Not configured")
			fmt.Println("   Run 'autocommit init' to configure")
		} else {
			if hasAPIKey(cfg.Provider) {
				fmt.Println("OK")

				fmt.Print("Provider connection... ")
				prov, err := provider.Get(cfg)
				if err != nil {
					fmt.Printf("Error: %s\n", err)
					issues++
				} else {
					if err := prov.Validate(); err != nil {
						fmt.Printf("Error: %s\n", err)
						issues++
					} else {
						fmt.Println("OK")
					}
				}
			} else {
				fmt.Println("API key not found")
				fmt.Printf("   Set %s environment variable\n", getAPIKeyEnvName(cfg.Provider))
				issues++
			}
		}
	}

	fmt.Print("Git hook... ")
	hookPath := ".git/hooks/prepare-commit-msg"
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		fmt.Println("Not installed (optional)")
	} else {
		content, _ := os.ReadFile(hookPath)
		if isOurHook(string(content)) {
			fmt.Println("Installed")
		} else {
			fmt.Println("Different hook installed")
		}
	}

	fmt.Println()
	if issues == 0 {
		fmt.Println("All checks passed!")
	} else {
		fmt.Printf("Found %d issue(s)\n", issues)
	}

	return nil
}

func hasAPIKey(providerName string) bool {
	return os.Getenv(getAPIKeyEnvName(providerName)) != ""
}

func getAPIKeyEnvName(providerName string) string {
	switch providerName {
	case "anthropic":
		return "ANTHROPIC_API_KEY"
	case "openai":
		return "OPENAI_API_KEY"
	case "gigachat":
		return "GIGACHAT_CLIENT_ID"
	case "yandexgpt":
		return "YANDEX_API_KEY"
	case "google":
		return "GOOGLE_API_KEY"
	case "mistral":
		return "MISTRAL_API_KEY"
	default:
		return "AUTOCOMMIT_API_KEY"
	}
}
