package command

import "context"

type Command interface {
	// GetFriendlyName returns the name to be used to call the program
	GetFriendlyName() string

	// GetDescription returns the description for the command
	GetDescription() string

	// Execute launches the command
	Execute(context.Context) error

	// defineFlags defines the flags used by the single command
	defineFlags()

	// checkFlagsAndDefaults checks if the flags are correctly set and assign default values
	checkFlagsAndDefaults() error
}

func GetRegisteredCommands() []Command {
	return []Command{
		&FindBranchCommand{},
		&AlignBranchCommand{},
	}
}
