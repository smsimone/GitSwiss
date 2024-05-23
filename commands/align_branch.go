package command

import (
	"context"
	"flag"
	"fmt"
	"log"

	"it.smaso/git_swiss/internal/git"
	"it.smaso/git_swiss/internal/utilities"
	"it.smaso/git_swiss/pool"
)

type AlignBranchCommand struct {
	Command
	source         *string
	target         *string
	directory      *string
	strategy       *string
	remote         *string
	parsedStrategy git.GitStrategy
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
	c.strategy = flag.String("strategy", "merge", "The strategy to use for alignment (merge | pull)")
	c.remote = flag.String("remote", "origin", "The remote to use for alignment")
}

func (c *AlignBranchCommand) CheckFlagsAndDefaults() error {
	if c.source == nil || len(*c.source) == 0 {
		return fmt.Errorf("missing required source branch")
	}

	if c.strategy != nil {
		strat, err := git.StrategyFromString(*c.strategy)
		if err != nil {
			return err
		}
		c.parsedStrategy = strat
	} else {
		c.parsedStrategy = git.MERGE_STRATEGY
	}

	if c.directory == nil {
		dir := "."
		c.directory = &dir
	}

	if c.remote == nil {
		defaultRemote := "origin"
		c.remote = &defaultRemote
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

	pool.Execute(
		func(path string) error {
			if err := git.Align(context.Background(), path, *c.source, *c.target, c.parsedStrategy, *c.remote); err != nil {
				log.Printf("Failed to align branch in %s: %s\n", path, err.Error())
				return err
			}
			return nil
		},
		*repositories,
	)

	return nil
}
