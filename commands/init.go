package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
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
			os.Exit(0)
		}

		// Get current directory name as default
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error: failed to get current directory: %v\n", err)
			os.Exit(1)
		}
		defaultName := filepath.Base(currentDir)

		// Initialize with defaults
		var projectName = defaultName
		var idPrefix = strings.ToUpper(defaultName)

		// Interactive form
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Project Name").
					Value(&projectName),

				huh.NewInput().
					Title("ID Prefix").
					Value(&idPrefix),
			),
		)

		err = form.Run()
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}

		// Create new project
		project := &Project{
			Name:         projectName,
			IDPrefix:     idPrefix,
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

		fmt.Printf("'%s' is ready for requirements\n", project.Name)
	},
}

// Project represents a collection of requirements for a Product Requirements Document
type Project struct {
	Name         string        `yaml:"name"`
	IDPrefix     string        `yaml:"id_prefix"`
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