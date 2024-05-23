package git

import (
	"fmt"

	"os"
	"os/exec"

	"it.smaso/git_swiss/internal/utilities"
)

type GitCommand struct {
	Path    string
	Options []string
}

func (c *GitCommand) getLoggableName() string {
	return c.Options[0]
}

func buildCommand(gitCmd GitCommand) *exec.Cmd {
	cmd := exec.Command("git", gitCmd.Options...)
	cmd.Dir = gitCmd.Path
	cmd.Env = os.Environ()

	return cmd
}

func Execute(gitCmd GitCommand) error {
	cmd := buildCommand(gitCmd)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to %s: %s", gitCmd.getLoggableName(), err.Error())
	}
	return nil
}

func ExecuteWithOutput(gitCmd GitCommand) (*string, error) {
	cmd := buildCommand(gitCmd)

	bytes, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	out := string(bytes)
	return &out, nil
}

// CheckGitRepo checks if the path given leads to a git repository
func CheckGitRepo(path string) error {
	if !utilities.ContainsFile(path, ".git") {
		return fmt.Errorf("not executing in a git repository")
	}
	return nil
}
