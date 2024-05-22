package command

import (
	"context"
	"flag"
	"fmt"

	"it.smaso/git_swiss/internal/git"
	"it.smaso/git_swiss/internal/utilities"
	"it.smaso/git_swiss/pool"
)

type CreateBranchCommand struct {
	Command
	directory *string
	source    *string
	target    *string
}

func (c *CreateBranchCommand) GetFriendlyName() string {
	return "create-branch"
}

// GetDescription returns the description for the command
func (c *CreateBranchCommand) GetDescription() string {
	return "creates a new branch named as required"
}

// defineFlags defines the flags used by the single command
func (c *CreateBranchCommand) DefineFlags() {
	c.directory = flag.String("directory", ".", "directory in which the command should be run (defaults to current directory)")
	c.source = flag.String("source", "", "branch from which checkout new branch")
	c.target = flag.String("target", "", "new branch name to be created")
}

// checkFlagsAndDefaults checks if the flags are correctly set and assign default values
func (c *CreateBranchCommand) CheckFlagsAndDefaults() error {
	if c.target == nil || len(*c.target) == 0 {
		return fmt.Errorf("target branch is required")
	}

	if c.directory == nil {
		dir := "."
		c.directory = &dir
	}

	return nil
}

// Execute launches the command
func (c *CreateBranchCommand) Execute(ctx context.Context) error {
	if err := StartupChecks(c); err != nil {
		return err
	}

	repositories, err := utilities.FindRepositories(context.Background(), *c.directory)
	if err != nil {
		return fmt.Errorf("failed to find directories: %s", err.Error())
	}

	pool.Execute(
		func(path string) error {
			return createBranch(path, c.source, *c.target)
		},
		*repositories,
	)

	return nil
}

func createBranch(path string, source *string, target string) error {
	var foundSource string
	if source != nil {
		res := git.BranchExists(context.Background(), path, *source)
		if res == nil {
			return fmt.Errorf("source branch does not exists")
		}
		foundSource = res.Name
	} else {
		current, err := git.CurrentBranch(context.Background(), path)
		if err != nil {
			return err
		}
		foundSource = *current
	}

	if err := git.Checkout(context.Background(), path, foundSource); err != nil {
		return fmt.Errorf("failed to checkout branch: %s", err.Error())
	}

	if err := git.Pull(context.Background(), path); err != nil {
		return fmt.Errorf("failed to pull project: %s", err.Error())
	}

	if err := git.CheckoutCreate(context.Background(), path, target); err != nil {
		return fmt.Errorf("failed to create new branch: %s", err.Error())
	}
	return nil

}
