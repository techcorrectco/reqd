package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new requirements project",
	Long:  `Initialize a new requirements project by creating a requirements.yaml file in the current directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		const filename = "requirements.yaml"

		// Check if file already exists
		if _, err := os.Stat(filename); err == nil {
			fmt.Printf("Error: %s already exists in this directory\n", filename)
			os.Exit(1)
		}

		// Get current directory name
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error: failed to get current directory: %v\n", err)
			os.Exit(1)
		}

		// Create new project
		project := &Project{
			Name:         filepath.Base(currentDir),
			Requirements: []Requirement{},
		}

		// Marshal to YAML
		data, err := yaml.Marshal(project)
		if err != nil {
			fmt.Printf("Error: failed to marshal project: %v\n", err)
			os.Exit(1)
		}

		// Write to file
		if err := os.WriteFile(filename, data, 0644); err != nil {
			fmt.Printf("Error: failed to write %s: %v\n", filename, err)
			os.Exit(1)
		}

		fmt.Printf("Initialized new requirements project '%s' in %s\n", project.Name, filename)
	},
}

// Project represents a collection of requirements for a Product Requirements Document
type Project struct {
	Name         string        `yaml:"name"`
	Requirements []Requirement `yaml:"requirements,omitempty"`
}

// Requirement represents a single requirement in a Product Requirements Document
type Requirement struct {
	ID          string        `yaml:"id"`
	Title       string        `yaml:"title"`
	Keyword     string        `yaml:"keyword"`
	Description string        `yaml:"description"`
	Children    []Requirement `yaml:"children,omitempty"`
}