package git

import (
	"context"
	"fmt"
	"strings"

	"it.smaso/git_swiss/internal/utilities"
)

// CurrentBranch returns the current branch of the repository
func CurrentBranch(ctx context.Context, path string) (*string, error) {
	if !utilities.ContainsFile(path, ".git") {
		return nil, fmt.Errorf("not executing in a git repository")
	}

	output, err := ExecuteWithOutput(GitCommand{
		Path:    path,
		Options: []string{"branch"},
	})

	if err != nil {
		return nil, fmt.Errorf("failed to exec git command: %s", err.Error())
	}

	lines := strings.Split(*output, "\n")
	for _, x := range lines {
		if strings.Contains(x, "*") {
			clean := strings.Trim(strings.Replace(x, "*", "", -1), " ")
			return &clean, nil
		}
	}

	return nil, fmt.Errorf("no branch found")
}
