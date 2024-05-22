package command

import (
	"context"
	"flag"
	"fmt"
	"os"

	"it.smaso/git_swiss/internal/git"
	"it.smaso/git_swiss/internal/utilities"
	"it.smaso/git_swiss/pool"
)

type AlignBranchCommand struct {
	Command
	source    *string
	target    *string
	directory *string
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
}

func (c *AlignBranchCommand) CheckFlagsAndDefaults() error {
	if c.source == nil || len(*c.source) == 0 {
		return fmt.Errorf("missing required target branch")
	}
	if c.target == nil || len(*c.target) == 0 {
		c.target = c.source
	}
	if c.directory == nil {
		dir := "."
		c.directory = &dir
	}

	return nil
}

func (c *AlignBranchCommand) Execute(ctx context.Context) error {
	if err := StartupChecks(c); err != nil {
		return err
	}

	repositories, err := utilities.FindRepositories(context.Background(), *c.directory)
	if err != nil {
		return fmt.Errorf("failed to find directories: %s", err.Error())
	}

	paths := []string{}
	for _, dir := range *repositories {
		if dir.Name() == *c.directory {
			paths = append(paths, *c.directory)
			continue
		}
		path := fmt.Sprintf("%s%s%s", *c.directory, string(os.PathSeparator), dir.Name())
		paths = append(paths, path)
	}

	pool.Execute(
		func(path string) error {
			if err := git.Align(context.Background(), path, *c.source, *c.target); err != nil {
				fmt.Printf("Failed to align branch in %s: %s\n", path, err.Error())
				return err
			}
			return nil
		},
		paths,
	)

	return nil
}
