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

type UpdateBranchCommand struct {
	Command
	directory *string
	branch    *string
	checkout  *bool
}

func (c *UpdateBranchCommand) GetFriendlyName() string {
	return "update-branch"
}

func (c *UpdateBranchCommand) GetDescription() string {
	return "Updates the given branch for each project in the given directory"
}

func (c *UpdateBranchCommand) DefineFlags() {
	c.directory = flag.String("directory", ".", "The project directory to update (defaults to the current directory)")
	c.branch = flag.String("branch", "", "The branch to update (defaults to the current branch)")
	c.checkout = flag.Bool("checkout", false, "If the repository is not on the selected branch, force checkout to the selected branch")
}

func (c *UpdateBranchCommand) CheckFlagsAndDefaults() error {
	if c.directory == nil {
		dir := "."
		c.directory = &dir
	}
	return nil
}

func (c *UpdateBranchCommand) Execute(ctx context.Context) error {
	if err := StartupChecks(c); err != nil {
		return err
	}

	repositories, err := utilities.FindRepositories(context.Background(), *c.directory)
	if err != nil {
		return err
	}

	errors := pool.Execute(
		func(path string) error {
			files, err := git.UncommittedFiles(context.Background(), path)
			if err != nil {
				log.Printf("failed to get pending changes for %s: %s\n", path, err.Error())
				return err
			}
			if len(files) > 0 {
				log.Printf("skipping %s, %d pending changes\n", path, len(files))
				return fmt.Errorf("there were some pending changes")
			}

			branch, err := git.CurrentBranch(context.Background(), path)
			if err != nil {
				log.Printf("failed to get current branch of %s: %s\n", path, err.Error())
				return err
			}

			if *branch != *c.branch {
				if *c.checkout {
					if err := git.Checkout(context.Background(), path, *c.branch); err != nil {
						log.Printf("failed checkout for %s: %s\n", path, err.Error())
						return err
					}
				} else {
					log.Printf("skipping %s, current branch is %s, expected %s\n", path, *branch, *c.branch)
					log.Printf("to force checkout, use '-checkout' flag")
					return fmt.Errorf("branch mismatch")
				}
			}

			if err := git.Pull(context.Background(), path); err != nil {
				log.Printf("failed to pull data %s: %s\n", path, err.Error())
				return err
			}

			log.Printf("Updated %s", path)

			return nil
		},
		*repositories,
	)

	for _, err := range errors {
		if err != nil {
			return err
		}
	}

	return nil
}
