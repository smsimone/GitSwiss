package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// UncommittedFiles returns the list of uncommitted files in the repository
func UncommittedFiles(ctx context.Context, path string) ([]string, error) {
	cmd := exec.Command("git", "status", "-s")
	cmd.Dir = path
	cmd.Env = os.Environ()

	bytes, err := cmd.Output()
	if err != nil {
		return []string{}, fmt.Errorf("failed to run git status: %s", err.Error())
	}
	content := strings.Split(string(bytes), "\n")
	return content, nil
}

// PendingChanges returns true if there are uncommitted files in the repository
func PendingChanges(ctx context.Context, path string) (bool, error) {
	files, err := UncommittedFiles(ctx, path)
	return len(files) > 0, err
}
