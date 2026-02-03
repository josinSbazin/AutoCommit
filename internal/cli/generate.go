package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/josinSbazin/AutoCommit/internal/config"
	"github.com/josinSbazin/AutoCommit/internal/git"
	"github.com/josinSbazin/AutoCommit/internal/prompt"
	"github.com/josinSbazin/AutoCommit/internal/provider"
	"github.com/josinSbazin/AutoCommit/internal/ui"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a commit message",
	RunE:  runGenerate,
}

func init() {
	generateCmd.Flags().BoolP("dry-run", "d", false, "Don't commit, just show message")
	generateCmd.Flags().StringP("output", "o", "", "Write message to file")
	generateCmd.Flags().Bool("hook-mode", false, "Run in git hook mode")
}

func runGenerate(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	if p, _ := cmd.Flags().GetString("provider"); p != "" {
		cfg.Provider = p
	}
	if m, _ := cmd.Flags().GetString("model"); m != "" {
		cfg.Model = m
	}

	if !git.IsRepo() {
		return fmt.Errorf("not a git repository")
	}

	diff, err := git.GetStagedDiff()
	if err != nil {
		return fmt.Errorf("failed to get diff: %w", err)
	}

	if diff.IsEmpty() {
		return fmt.Errorf("no staged changes. Use 'git add' first")
	}

	history, _ := git.GetCommitHistory(cfg.Context.HistoryCount)
	branch, _ := git.GetCurrentBranch()

	promptBuilder := prompt.NewBuilder(cfg)
	promptText := promptBuilder.Build(diff, history, branch)

	prov, err := provider.Get(cfg)
	if err != nil {
		return fmt.Errorf("failed to get provider: %w", err)
	}

	spinner := ui.NewSpinner("Generating commit message...")
	spinner.Start()

	message, err := prov.Generate(ctx, promptText)
	spinner.Stop()

	if err != nil {
		return fmt.Errorf("failed to generate message: %w", err)
	}

	dryRun, _ := cmd.Flags().GetBool("dry-run")
	outputFile, _ := cmd.Flags().GetString("output")
	hookMode, _ := cmd.Flags().GetBool("hook-mode")

	if hookMode && outputFile != "" {
		return os.WriteFile(outputFile, []byte(message), 0644)
	}

	if outputFile != "" {
		if err := os.WriteFile(outputFile, []byte(message), 0644); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}
		fmt.Printf("Message written to %s\n", outputFile)
		return nil
	}

	if dryRun {
		fmt.Println(message)
		return nil
	}

	return runInteractive(ctx, cfg, prov, promptText, message)
}

func runInteractive(ctx context.Context, cfg *config.Config, prov provider.Provider, promptText, message string) error {
	for {
		ui.PrintCommitMessage(message)
		action := ui.AskAction()

		switch action {
		case ui.ActionAccept:
			if err := git.CreateCommit(message); err != nil {
				return fmt.Errorf("failed to commit: %w", err)
			}
			ui.PrintSuccess("Committed!")
			return nil

		case ui.ActionEdit:
			edited, err := ui.EditInEditor(message)
			if err != nil {
				return fmt.Errorf("failed to edit: %w", err)
			}
			message = edited

		case ui.ActionRegenerate:
			spinner := ui.NewSpinner("Regenerating...")
			spinner.Start()
			newMessage, err := prov.Generate(ctx, promptText)
			spinner.Stop()
			if err != nil {
				ui.PrintError("Failed to regenerate: " + err.Error())
				continue
			}
			message = newMessage

		case ui.ActionQuit:
			fmt.Println("Aborted.")
			return nil
		}
	}
}
