package git

import (
	"context"
	"fmt"
	"strings"
)

// UncommittedFiles returns the list of uncommitted files in the repository
func UncommittedFiles(ctx context.Context, path string) ([]string, error) {
	output, err := ExecuteWithOutput(GitCommand{
		Path:    path,
		Options: []string{"status", "-s"},
	})

	if err != nil {
		return []string{}, fmt.Errorf("failed to run git status: %s", err.Error())
	}

	content := strings.Split(*output, "\n")

	files := []string{}
	for _, x := range content {
		if len(strings.TrimSpace(x)) > 0 {
			files = append(files, x)
		}
	}
	return files, nil
}

// PendingChanges returns true if there are uncommitted files in the repository
func PendingChanges(ctx context.Context, path string) (bool, error) {
	files, err := UncommittedFiles(ctx, path)
	return len(files) > 0, err
}
