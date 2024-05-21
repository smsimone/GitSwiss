package git

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

type BranchResult struct {
	Name      string
	IsCurrent bool
}

// BranchExists checks if the required branch exists in the projectDirectory
// If the branch exists, it returns the full branch name and true
func BranchExists(ctx context.Context, projectDirectory, branchName string) *BranchResult {
	cmd := exec.Command("git", "branch", "-a")
	cmd.Dir = projectDirectory
	cmd.Env = os.Environ()

	bytes, err := cmd.Output()
	if err != nil {
		fmt.Printf("Failed to run command: %s\n", err.Error())
		return nil
	}
	outLines := strings.Split(string(bytes), "\n")
	for _, line := range outLines {
		if strings.Contains(line, branchName) {
			isCurrent := strings.Contains(line, "*")
			line = strings.Trim(strings.Replace(line, "*", "", -1), " ")
			return &BranchResult{Name: line, IsCurrent: isCurrent}
		}
	}

	return nil
}
