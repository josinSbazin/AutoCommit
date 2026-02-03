package prompt

import (
	"fmt"
	"strings"

	"github.com/josinSbazin/AutoCommit/internal/config"
	"github.com/josinSbazin/AutoCommit/internal/git"
)

type Builder struct {
	cfg *config.Config
}

func NewBuilder(cfg *config.Config) *Builder {
	return &Builder{cfg: cfg}
}

// Build constructs the prompt from diff and context
func (b *Builder) Build(diff *git.DiffResult, history []git.Commit, branch string) string {
	var sb strings.Builder

	// System instructions
	sb.WriteString("You are an expert at writing clear, concise git commit messages.\n\n")

	// Task
	sb.WriteString("## Task\n")
	sb.WriteString("Analyze the git diff below and generate a commit message.\n\n")

	// Style instructions
	sb.WriteString("## Style\n")
	sb.WriteString(b.getStyleInstructions())
	sb.WriteString("\n\n")

	// Language
	sb.WriteString("## Language\n")
	sb.WriteString(b.getLanguageInstructions())
	sb.WriteString("\n\n")

	// Context
	if b.cfg.Context.IncludeHistory && len(history) > 0 {
		sb.WriteString("## Project Commit Style (for reference)\n")
		sb.WriteString("Recent commits in this repository:\n")
		for i, c := range history {
			if i >= 5 {
				break
			}
			sb.WriteString(fmt.Sprintf("- %s\n", c.Subject))
		}
		sb.WriteString("\n")
	}

	if b.cfg.Context.IncludeBranch && branch != "" {
		sb.WriteString("## Branch\n")
		sb.WriteString(fmt.Sprintf("Current branch: %s\n\n", branch))
	}

	// Diff stats
	if b.cfg.Context.IncludeDiffStats {
		sb.WriteString("## Changes Summary\n")
		sb.WriteString(diff.Summary())
		sb.WriteString("\n")
	}

	// Full diff
	sb.WriteString("## Git Diff\n")
	sb.WriteString("```diff\n")

	// Truncate if too long
	rawDiff := diff.RawDiff
	if len(rawDiff) > 10000 {
		rawDiff = rawDiff[:10000] + "\n... (truncated)"
	}
	sb.WriteString(rawDiff)
	sb.WriteString("\n```\n\n")

	// Custom instructions
	if b.cfg.Instructions != "" {
		sb.WriteString("## Additional Instructions\n")
		sb.WriteString(b.cfg.Instructions)
		sb.WriteString("\n\n")
	}

	// Output format
	sb.WriteString("## Output Format\n")
	sb.WriteString(b.getOutputInstructions())

	return sb.String()
}

func (b *Builder) getStyleInstructions() string {
	switch b.cfg.Style {
	case "conventional":
		return `Use Conventional Commits format:
- Format: type(scope): description
- Types: feat, fix, docs, style, refactor, test, chore, perf, ci, build
- Scope is optional but recommended
- Description should be imperative mood ("add" not "added")
- First letter lowercase
- No period at the end`

	case "simple":
		return `Use simple format:
- Short, descriptive message
- Imperative mood ("add" not "added")
- Max 72 characters
- No need for type prefixes`

	case "detailed":
		return `Use detailed format:
- First line: summary (max 72 chars)
- Empty line
- Body: bullet points explaining what and why
- Be thorough but concise`

	default:
		return `Write a clear, concise commit message.
- Imperative mood
- Max 72 characters for first line`
	}
}

func (b *Builder) getLanguageInstructions() string {
	switch b.cfg.Language {
	case "ru":
		return "Write the commit message in Russian (Русский)"
	case "de":
		return "Write the commit message in German (Deutsch)"
	case "fr":
		return "Write the commit message in French (Français)"
	case "es":
		return "Write the commit message in Spanish (Español)"
	case "zh":
		return "Write the commit message in Chinese (中文)"
	case "ja":
		return "Write the commit message in Japanese (日本語)"
	default:
		return "Write the commit message in English"
	}
}

func (b *Builder) getOutputInstructions() string {
	var sb strings.Builder

	sb.WriteString("Return ONLY the commit message, nothing else.\n")
	sb.WriteString("No explanations, no markdown formatting, no quotes.\n\n")

	sb.WriteString(fmt.Sprintf("First line: subject (max %d characters)\n", b.cfg.MaxSubjectLength))

	if b.cfg.IncludeBody {
		sb.WriteString("Then empty line, then body with bullet points.\n")
		sb.WriteString(fmt.Sprintf("Body should be max %d characters total.\n", b.cfg.MaxBodyLength))
	} else {
		sb.WriteString("Only the subject line, no body.\n")
	}

	return sb.String()
}
