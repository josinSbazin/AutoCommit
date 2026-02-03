package git

import (
	"os/exec"
	"strconv"
	"strings"
)

// Commit represents a git commit
type Commit struct {
	Hash    string
	Subject string
	Body    string
	Author  string
}

func GetCommitHistory(n int) ([]Commit, error) {
	if n <= 0 {
		n = 10
	}

	format := "%h|%s|%b|%an"
	cmd := exec.Command("git", "log", "-n", strconv.Itoa(n), "--format="+format, "--no-merges")

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
