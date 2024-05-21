package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"it.smaso/git_swiss/internal/utilities"
)

func Pull(ctx context.Context, path string) error {
	if !utilities.ContainsFile(path, ".git") {
		return fmt.Errorf("not executing in a git repository")
	}

	cmd := exec.Command("git", "pull")
	cmd.Dir = path
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull: %s", err.Error())
	}

	return nil
}
