package utilities

import (
	"fmt"
	"os"
)

func ContainsFile(path, name string) bool {
	if len(name) == 0 {
		panic("name cannot be empty")
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Printf("failed to read directory '%s': %s\n", path, err.Error())
		return false
	}

	for _, entry := range entries {
		if entry.Name() == name {
			return true
		}
	}
	return false
}
