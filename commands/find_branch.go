package command

import (
	"context"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"strings"

	"it.smaso/git_utilities/internal/git"
	"it.smaso/git_utilities/internal/utilities"
	"it.smaso/git_utilities/pool"
)

type branchRes struct {
	project string
	branch  git.BranchResult
}

type FindBranchCommand struct {
	Command
	requestedBranch *string
	directory       *string
	helpMsg         *bool
}

func (cmd *FindBranchCommand) defineFlags() {
	cmd.requestedBranch = flag.String("branch", "", "The branch to check")
	cmd.directory = flag.String("directory", "", "The directory to check the branch in")
	cmd.helpMsg = flag.Bool("help", false, "Print the help message for the command")
}

func (cmd *FindBranchCommand) checkFlagsAndDefaults() error {
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
	return "find_branch"
}

func (cmd *FindBranchCommand) GetDescription() string {
	return "Scans all the project in the given folder and returns all the projects that contains the branch"
}

func (cmd *FindBranchCommand) Execute(ctx context.Context) error {
	cmd.defineFlags()
	flag.CommandLine.Parse(os.Args[2:])
	if err := cmd.checkFlagsAndDefaults(); err != nil {
		fmt.Println("Failed to parse flags\n", err.Error())
		flag.PrintDefaults()
		return nil
	}

	if cmd.helpMsg != nil && *cmd.helpMsg {
		flag.PrintDefaults()
		return nil
	}

	folders, err := utilities.FindFolders(context.Background(), *cmd.directory)
	if err != nil {
		fmt.Printf("Failed to read folder '%s': %s\n", *cmd.directory, err.Error())
		return err
	}

	compos := strings.Split(*cmd.directory, string(os.PathSeparator))
	dirname := compos[len(compos)-1]

	names := pool.Execute(
		func(a fs.DirEntry) *branchRes {
			path := fmt.Sprintf("%s/%s", *cmd.directory, a.Name())
			if strings.HasSuffix(path, dirname) {
				path = *cmd.directory
			}
			if !utilities.ContainsFile(path, ".git") {
				return nil
			}
			res := git.BranchExists(
				context.Background(),
				path,
				*cmd.requestedBranch,
			)
			if res == nil {
				return nil
			}
			return &branchRes{project: a.Name(), branch: *res}
		},
		*folders,
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
