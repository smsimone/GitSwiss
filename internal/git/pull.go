package git

import (
	"context"
	"fmt"
	"os/exec"

	"it.smaso/git_utilities/configs"
	"it.smaso/git_utilities/internal/utilities"
)

func Pull(ctx context.Context, path string) error {
	if !utilities.ContainsFile(path, ".git") {
		return fmt.Errorf("not executing in a git repository")
	}

	cmd := exec.Command(configs.Instance().GitExec, "pull")
	cmd.Dir = path

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull: %s", err.Error())
	}

	return nil
}
