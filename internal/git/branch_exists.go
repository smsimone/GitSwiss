package git

import (
	"context"
	"fmt"
	"strings"
)

type BranchResult struct {
	Name      string
	IsCurrent bool
}

// BranchExists checks if the required branch exists in the projectDirectory
// If the branch exists, it returns the full branch name and true
func BranchExists(ctx context.Context, projectDirectory, branchName string) *BranchResult {
	output, err := ExecuteWithOutput(GitCommand{
		Path:    projectDirectory,
		Options: []string{"branch", "-a"},
	})

	if err != nil {
		fmt.Printf("Failed to run command: %s\n", err.Error())
		return nil
	}

	outLines := strings.Split(*output, "\n")
	for _, line := range outLines {
		if strings.Contains(line, branchName) {
			isCurrent := strings.Contains(line, "*")
			line = strings.Trim(strings.Replace(line, "*", "", -1), " ")
			return &BranchResult{Name: line, IsCurrent: isCurrent}
		}
	}

	return nil
}
