package git

import (
	"os/exec"
	"strconv"
	"strings"
)

// DiffResult holds parsed git diff information
type DiffResult struct {
	Files    []FileChange
	Stats    DiffStats
	RawDiff  string
	IsBinary bool
}

// FileChange represents changes to a single file
type FileChange struct {
	Path       string
	OldPath    string // for renames
	Status     string // added, modified, deleted, renamed
	Additions  int
	Deletions  int
	IsBinary   bool
	DiffChunk  string
}

// DiffStats holds overall diff statistics
type DiffStats struct {
	FilesChanged int
	Additions    int
	Deletions    int
}

// IsEmpty returns true if there are no changes
func (d *DiffResult) IsEmpty() bool {
	return len(d.Files) == 0
}

// GetStagedDiff returns the diff of staged changes
func GetStagedDiff() (*DiffResult, error) {
	// Get raw diff
	cmd := exec.Command("git", "diff", "--cached", "--unified=3")
	rawDiff, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// Get file stats
	cmd = exec.Command("git", "diff", "--cached", "--numstat")
	numstat, _ := cmd.Output()

	// Get file status
	cmd = exec.Command("git", "diff", "--cached", "--name-status")
	nameStatus, _ := cmd.Output()

	result := &DiffResult{
		RawDiff: string(rawDiff),
		Files:   parseNameStatus(string(nameStatus)),
	}

	// Parse numstat and merge with files
	parseNumstat(string(numstat), result)

	// Calculate totals
	for _, f := range result.Files {
		result.Stats.Additions += f.Additions
		result.Stats.Deletions += f.Deletions
		if f.IsBinary {
			result.IsBinary = true
		}
	}
	result.Stats.FilesChanged = len(result.Files)

	return result, nil
}

// parseNameStatus parses git diff --name-status output
func parseNameStatus(output string) []FileChange {
	var files []FileChange
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		fc := FileChange{
			Path: parts[1],
		}

		switch parts[0][0] {
		case 'A':
			fc.Status = "added"
		case 'M':
			fc.Status = "modified"
		case 'D':
			fc.Status = "deleted"
		case 'R':
			fc.Status = "renamed"
			if len(parts) >= 3 {
				fc.OldPath = parts[1]
				fc.Path = parts[2]
			}
		case 'C':
			fc.Status = "copied"
		default:
			fc.Status = "modified"
		}

		files = append(files, fc)
	}

	return files
}

// parseNumstat parses git diff --numstat output and updates files
func parseNumstat(output string, result *DiffResult) {
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		path := parts[2]

		// Handle renames (path is "old => new")
		if strings.Contains(path, " => ") {
			pathParts := strings.Split(path, " => ")
			if len(pathParts) >= 2 {
				path = pathParts[1]
			}
		}

		// Find matching file
		for i := range result.Files {
			if result.Files[i].Path == path {
				// Binary files show "-" for additions/deletions
				if parts[0] == "-" {
					result.Files[i].IsBinary = true
				} else {
					result.Files[i].Additions, _ = strconv.Atoi(parts[0])
					result.Files[i].Deletions, _ = strconv.Atoi(parts[1])
				}
				break
			}
		}
	}
}

// Summary returns a brief summary of the diff
func (d *DiffResult) Summary() string {
	var sb strings.Builder

	sb.WriteString("Changed files:\n")
	for _, f := range d.Files {
		sb.WriteString("  ")
		switch f.Status {
		case "added":
			sb.WriteString("[+] ")
		case "deleted":
			sb.WriteString("[-] ")
		case "modified":
			sb.WriteString("[~] ")
		case "renamed":
			sb.WriteString("[R] ")
		}
		sb.WriteString(f.Path)
		if f.Additions > 0 || f.Deletions > 0 {
			sb.WriteString(" (")
			if f.Additions > 0 {
				sb.WriteString("+")
				sb.WriteString(strconv.Itoa(f.Additions))
			}
			if f.Deletions > 0 {
				if f.Additions > 0 {
					sb.WriteString("/")
				}
				sb.WriteString("-")
				sb.WriteString(strconv.Itoa(f.Deletions))
			}
			sb.WriteString(")")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}
