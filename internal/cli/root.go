package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	versionStr string
	commitStr  string
	dateStr    string
)

func SetVersionInfo(version, commit, date string) {
	versionStr = version
	commitStr = commit
	dateStr = date
}

var rootCmd = &cobra.Command{
	Use:   "autocommit",
	Short: "AI-powered git commit messages",
	Long: `AutoCommit generates commit messages using LLMs.

Examples:
  autocommit              # Generate and commit interactively
  autocommit generate     # Just show the message
  autocommit init         # Setup for current project`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runGenerate(cmd, args)
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("provider", "p", "", "LLM provider")
	rootCmd.PersistentFlags().StringP("model", "m", "", "Model name")
	rootCmd.PersistentFlags().StringP("config", "c", "", "Config file path")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose output")

	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(hookCmd)
	rootCmd.AddCommand(doctorCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("autocommit %s\n", versionStr)
		fmt.Printf("  commit: %s\n", commitStr)
		fmt.Printf("  built:  %s\n", dateStr)
	},
}
