package main

import (
	"fmt"
	"os"

	"github.com/techcorrectco/reqd/commands"
	"github.com/techcorrectco/reqd/internal/types"
	"gopkg.in/yaml.v3"
)

func main() {
	// Try to load existing project if requirements.yaml exists
	project := loadProjectIfExists()

	// Store project reference for commands to use
	_ = project

	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func loadProjectIfExists() *types.Project {
	const filename = "requirements.yaml"

	// Check if requirements.yaml exists
	if _, err := os.Stat(filename); err == nil {
		// File exists, unmarshal it
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to read %s: %v\n", filename, err)
			return nil
		}

		var project types.Project
		if err := yaml.Unmarshal(data, &project); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to unmarshal %s: %v\n", filename, err)
			return nil
		}

		return &project
	}

	// File doesn't exist, return nil
	return nil
}