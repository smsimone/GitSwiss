package command

import (
	"context"
	"flag"
	"fmt"

	"it.smaso/git_swiss/internal/git"
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

	return nil
}

// Execute launches the command
func (c *CreateBranchCommand) Execute(ctx context.Context) error {
	if err := StartupChecks(c); err != nil {
		return err
	}

	var dir string
	if c.directory != nil {
		dir = *c.directory
	} else {
		dir = "."
	}

	var source string
	if c.source != nil {
		res := git.BranchExists(context.Background(), dir, *c.source)
		if res == nil {
			return fmt.Errorf("source branch does not exists")
		}
		source = res.Name
	} else {
		current, err := git.CurrentBranch(context.Background(), dir)
		if err != nil {
			return err
		}
		source = *current
	}

	if err := git.Checkout(context.Background(), dir, source); err != nil {
		return fmt.Errorf("failed to checkout branch: %s", err.Error())
	}

	if err := git.Pull(context.Background(), dir); err != nil {
		return fmt.Errorf("failed to pull project: %s", err.Error())
	}

	if err := git.CheckoutCreate(context.Background(), dir, *c.target); err != nil {
		return fmt.Errorf("failed to create new branch: %s", err.Error())
	}

	return nil
}
