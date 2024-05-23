package command

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"it.smaso/git_swiss/internal/git"
	"it.smaso/git_swiss/internal/utilities"
	"it.smaso/git_swiss/pool"
)

type branchRes struct {
	project string
	branch  git.BranchResult
}

type FindBranchCommand struct {
	Command
	requestedBranch *string
	directory       *string
	checkout        *bool
}

func (cmd *FindBranchCommand) DefineFlags() {
	cmd.requestedBranch = flag.String("branch", "", "The branch to check")
	cmd.directory = flag.String("directory", "", "The directory to check the branch in")
	cmd.checkout = flag.Bool("checkout", false, "If specified, the branch will be checked out on each project found")
}

func (cmd *FindBranchCommand) CheckFlagsAndDefaults() error {
	if cmd.requestedBranch == nil || len(*cmd.requestedBranch) == 0 {
		return fmt.Errorf("Branch name is required")
	}

	if cmd.directory == nil || len(*cmd.directory) == 0 {
		dir := "."
		cmd.directory = &dir
	}
	return nil
}

func (cmd *FindBranchCommand) GetFriendlyName() string {
	return "find-branch"
}

func (cmd *FindBranchCommand) GetDescription() string {
	return "Scans all the project in the given folder and returns all the projects that contains the branch"
}

func (cmd *FindBranchCommand) Execute(ctx context.Context) error {
	if err := StartupChecks(cmd); err != nil {
		return err
	}

	repositories, err := utilities.FindRepositories(context.Background(), *cmd.directory)
	if err != nil {
		fmt.Printf("Failed to find repositories '%s': %s\n", *cmd.directory, err.Error())
		return err
	}

	names := pool.Execute(
		func(path string) *branchRes {
			res := git.BranchExists(
				context.Background(),
				path,
				*cmd.requestedBranch,
			)
			if res == nil {
				return nil
			}
			comps := strings.Split(path, string(os.PathSeparator))
			name := comps[len(comps)-1]

			if *cmd.checkout {
				err := git.Checkout(context.Background(), path, res.Name)
				if err != nil {
					fmt.Printf("Failed to checkout branch '%s' in '%s': %s\n", res.Name, path, err.Error())
					return nil
				}
			}

			return &branchRes{project: name, branch: *res}
		},
		*repositories,
	)

	for _, x := range names {
		if x != nil {
			sep := "."
			signal := "*"
			if !x.branch.IsCurrent {
				signal = sep
			}
			fmt.Printf("%s%s%s%s\n", x.project, strings.Repeat(sep, 40-len(x.project)), signal, x.branch.Name)
		}
	}
	return nil
}
