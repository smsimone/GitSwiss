package command

import (
	"context"
	"flag"
	"fmt"
	"os"
)

type Command interface {
	// GetFriendlyName returns the name to be used to call the program
	GetFriendlyName() string

	// GetDescription returns the description for the command
	GetDescription() string

	// Execute launches the command
	Execute(context.Context) error

	// defineFlags defines the flags used by the single command
	DefineFlags()

	// checkFlagsAndDefaults checks if the flags are correctly set and assign default values
	CheckFlagsAndDefaults() error
}

// StartupChecks runs pre-exec checks which are common between all commands
func StartupChecks(cmd Command) error {
	cmd.DefineFlags()
	flag.CommandLine.Parse(os.Args[2:])
	if err := cmd.CheckFlagsAndDefaults(); err != nil {
		fmt.Println("Failed to parse flags\n", err.Error())
		flag.PrintDefaults()
		return err
	}
	return nil
}

func GetRegisteredCommands() []Command {
	return []Command{
		&FindBranchCommand{},
		&AlignBranchCommand{},
		&CreateBranchCommand{},
		&UpdateBranchCommand{},
	}
}
