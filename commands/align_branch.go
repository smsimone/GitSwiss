package command

import (
	"context"
	"flag"
	"fmt"

	"it.smaso/git_utilities/internal/git"
	"it.smaso/git_utilities/internal/utilities"
)

type AlignBranchCommand struct {
	Command
	source    *string
	target    *string
	directory *string
	helpMsg   *bool
}

func (c *AlignBranchCommand) GetFriendlyName() string {
	return "align-branch"
}

func (c *AlignBranchCommand) GetDescription() string {
	return "Aligns the target branch with the source branch"
}

func (c *AlignBranchCommand) DefineFlags() {
	c.source = flag.String("source", "", "The branch to align from")
	c.target = flag.String("target", "", "The branch to align to (defaults to the current branch)")
	c.directory = flag.String("directory", ".", "The project directory to align (defaults to the current directory)")
	c.helpMsg = flag.Bool("help", false, "Print the help message for the command")
}

func (c *AlignBranchCommand) CheckFlagsAndDefaults() error {
	if c.source == nil || len(*c.source) == 0 {
		return fmt.Errorf("missing required source branch")
	}

	return nil
}

func (c *AlignBranchCommand) Execute(ctx context.Context) error {
	if err := StartupChecks(c); err != nil {
		return err
	}

	if c.helpMsg != nil && *c.helpMsg {
		flag.PrintDefaults()
		return nil
	}

	dirs, err := utilities.FindFolders(context.Background(), *c.directory)
	if err != nil {
		return fmt.Errorf("failed to find directories: %s", err.Error())
	}

	for _, dir := range *dirs {
		path := fmt.Sprintf("%s/%s", *c.directory, dir.Name())
		if !utilities.ContainsFile(path, ".git") {
			if err := git.Align(context.Background(), path, *c.source, *c.target); err != nil {
				fmt.Printf("Failed to align branch in %s: %s\n", dir, err.Error())
			}
		}
	}

	return git.Align(context.Background(), *c.directory, *c.source, *c.target)
}
