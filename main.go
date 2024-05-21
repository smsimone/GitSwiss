package main

import (
	"context"
	"fmt"
	"os"

	command "it.smaso/git_swiss/commands"
	"it.smaso/git_swiss/pool"
)

func main() {
	if len(os.Args) == 1 {
		printAvailableCommands()
		return
	}

	pool.NewPool(pool.WithMaxRunners(12))

	for _, cmd := range command.GetRegisteredCommands() {
		if cmd.GetFriendlyName() == os.Args[1] {
			ctx := context.Background()
			if err := cmd.Execute(ctx); err != nil {
				panic(err)
			}
			return
		}
	}
	fmt.Printf("Command %s is not registered\n", os.Args[1])
	printAvailableCommands()
}

func printAvailableCommands() {
	fmt.Println("Missing command")
	fmt.Println("Usage: go_utilities <program_name>")
	fmt.Println("")
	fmt.Println("Available programs are:")
	for _, cmd := range command.GetRegisteredCommands() {
		fmt.Printf("%s: %s\n", cmd.GetFriendlyName(), cmd.GetDescription())
	}
}
