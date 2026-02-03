package git

import (
	"os/exec"
	"strings"
)

// Commit represents a git commit
type Commit struct {
	Hash    string
	Subject string
	Body    string
	Author  string
}

// GetCommitHistory returns the last n commits
func GetCommitHistory(n int) ([]Commit, error) {
	if n <= 0 {
		n = 10
	}

	// Format: hash|subject|body|author
	format := "%h|%s|%b|%an"
	cmd := exec.Command("git", "log", "-n", string(rune('0'+n)), "--format="+format, "--no-merges")

	// For n > 9
	if n > 9 {
		cmd = exec.Command("git", "log", "-n", formatInt(n), "--format="+format, "--no-merges")
	}

	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var commits []Commit
	entries := strings.Split(strings.TrimSpace(string(out)), "\n")

	for _, entry := range entries {
		if entry == "" {
			continue
		}

		parts := strings.SplitN(entry, "|", 4)
		if len(parts) < 4 {
			continue
		}

		commits = append(commits, Commit{
			Hash:    parts[0],
			Subject: parts[1],
			Body:    parts[2],
			Author:  parts[3],
		})
	}

	return commits, nil
}

// GetCommitMessagesForStyle returns recent commit messages for style analysis
func GetCommitMessagesForStyle(n int) ([]string, error) {
	commits, err := GetCommitHistory(n)
	if err != nil {
		return nil, err
	}

	var messages []string
	for _, c := range commits {
		messages = append(messages, c.Subject)
	}
	return messages, nil
}

func formatInt(n int) string {
	if n < 10 {
		return string(rune('0' + n))
	}
	result := ""
	for n > 0 {
		result = string(rune('0'+n%10)) + result
		n /= 10
	}
	return result
}
