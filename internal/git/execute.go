package git

import (
	"fmt"
	"os"
	"os/exec"
)

func Execute(path string, commands ...string) error {
	cmd := exec.Command("git", commands...)
	cmd.Dir = path
	cmd.Env = os.Environ()

	if err := cmd.Run(); err != nil {
		firstCommand := commands[0]
		return fmt.Errorf("failed to %s: %s", firstCommand, err.Error())
	}
	return nil
}
