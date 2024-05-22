package git

import (
	"context"
	"fmt"

	"it.smaso/git_swiss/internal/utilities"
)

// If the provided path is not a git repository
// an error will be returned.
func CheckGitRepo(path string) error {
	if !utilities.ContainsFile(path, ".git") {
		return fmt.Errorf("not executing in a git repository")
	}
	return nil
}

func Pull(ctx context.Context, path string) error {
	err := CheckGitRepo(path)
	if err == nil {
		return Execute(path, "pull")
	}
	return err
}

func MergePull(ctx context.Context, path string, remote string, source string) error {
	err := CheckGitRepo(path)
	if err == nil {
		return Execute(path, "pull", remote, source)
	}
	return err
}
