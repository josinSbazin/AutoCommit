package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var hookCmd = &cobra.Command{
	Use:   "hook",
	Short: "Manage git hooks",
}

var hookInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install git hook",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := installHook(); err != nil {
			return err
		}
		fmt.Println("Git hook installed")
		return nil
	},
}

var hookUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall git hook",
	RunE: func(cmd *cobra.Command, args []string) error {
		hookPath := ".git/hooks/prepare-commit-msg"

		if _, err := os.Stat(hookPath); os.IsNotExist(err) {
			fmt.Println("No hook installed")
			return nil
		}

		content, err := os.ReadFile(hookPath)
		if err != nil {
			return err
		}

		if !isOurHook(string(content)) {
			return fmt.Errorf("hook exists but wasn't installed by autocommit")
		}

		if err := os.Remove(hookPath); err != nil {
			return err
		}

		fmt.Println("Git hook uninstalled")
		return nil
	},
}

var hookStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check hook status",
	RunE: func(cmd *cobra.Command, args []string) error {
		hookPath := ".git/hooks/prepare-commit-msg"

		if _, err := os.Stat(hookPath); os.IsNotExist(err) {
			fmt.Println("Hook not installed")
			return nil
		}

		content, err := os.ReadFile(hookPath)
		if err != nil {
			return err
		}

		if isOurHook(string(content)) {
			fmt.Println("AutoCommit hook is installed")
		} else {
			fmt.Println("Hook exists but is not managed by AutoCommit")
		}

		return nil
	},
}

func init() {
	hookCmd.AddCommand(hookInstallCmd)
	hookCmd.AddCommand(hookUninstallCmd)
	hookCmd.AddCommand(hookStatusCmd)
}

func isOurHook(content string) bool {
	return strings.Contains(content, "autocommit")
}
