package git

import (
	"context"
	"fmt"

	"it.smaso/git_swiss/internal/utilities"
)

// Merge merges the source branch into the current branch
func Merge(ctx context.Context, path, source string) error {
	if !utilities.ContainsFile(path, ".git") {
		return fmt.Errorf("not executing in a git repository")
	}

	err := Execute(GitCommand{
		Path:    path,
		Options: []string{"merge", source},
	})

	if err != nil {
		return fmt.Errorf("failed to merge branch: %s", err.Error())
	}

	return nil
}
