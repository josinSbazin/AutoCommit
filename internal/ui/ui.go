package ui

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// Action represents user action in interactive mode
type Action int

const (
	ActionAccept Action = iota
	ActionEdit
	ActionRegenerate
	ActionQuit
)

// Colors for terminal output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorCyan   = "\033[36m"
	colorGray   = "\033[90m"
	colorBold   = "\033[1m"
)

// PrintCommitMessage displays the generated commit message
func PrintCommitMessage(message string) {
	fmt.Println()
	fmt.Println(colorCyan + "Generated commit message:" + colorReset)
	fmt.Println(colorGray + "─────────────────────────────────────" + colorReset)
	fmt.Println(message)
	fmt.Println(colorGray + "─────────────────────────────────────" + colorReset)
	fmt.Println()
}

// AskAction prompts user for action
func AskAction() Action {
	fmt.Print(colorBold + "[Enter]" + colorReset + " Accept  ")
	fmt.Print(colorBold + "[e]" + colorReset + " Edit  ")
	fmt.Print(colorBold + "[r]" + colorReset + " Regenerate  ")
	fmt.Print(colorBold + "[q]" + colorReset + " Quit")
	fmt.Print("\n> ")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	switch input {
	case "", "y", "yes":
		return ActionAccept
	case "e", "edit":
		return ActionEdit
	case "r", "regenerate", "retry":
		return ActionRegenerate
	case "q", "quit", "exit", "n", "no":
		return ActionQuit
	default:
		return ActionAccept
	}
}

// EditInEditor opens the message in the default editor
func EditInEditor(message string) (string, error) {
	// Create temp file
	tmpFile, err := os.CreateTemp("", "autocommit-*.txt")
	if err != nil {
		return "", err
	}
	defer os.Remove(tmpFile.Name())

	// Write message
	if _, err := tmpFile.WriteString(message); err != nil {
		return "", err
	}
	tmpFile.Close()

	// Get editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = os.Getenv("VISUAL")
	}
	if editor == "" {
		// Try common editors
		editors := []string{"code", "vim", "nano", "notepad"}
		for _, e := range editors {
			if _, err := exec.LookPath(e); err == nil {
				editor = e
				break
			}
		}
	}
	if editor == "" {
		editor = "notepad" // Windows fallback
	}

	// Open editor
	cmd := exec.Command(editor, tmpFile.Name())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Read edited content
	content, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(content)), nil
}

// PrintSuccess prints a success message
func PrintSuccess(msg string) {
	fmt.Println(colorGreen + "✓ " + msg + colorReset)
}

// PrintError prints an error message
func PrintError(msg string) {
	fmt.Println(colorRed + "✗ " + msg + colorReset)
}

// PrintWarning prints a warning message
func PrintWarning(msg string) {
	fmt.Println(colorYellow + "⚠ " + msg + colorReset)
}

// Confirm asks for yes/no confirmation
func Confirm(question string) bool {
	fmt.Printf("%s [y/N]: ", question)

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))

	return input == "y" || input == "yes"
}
