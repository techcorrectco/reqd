package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/techcorrectco/reqd/internal/types"
	"gopkg.in/yaml.v3"
)

var InitCmd = &cobra.Command{
	Use:     "init",
	Aliases: []string{"i"},
	Short:   "Initialize a new requirements project",
	Long:    `Initialize a new requirements project by creating a requirements.yaml file in the current directory.`,
	Run: func(cmd *cobra.Command, args []string) {
		const filename = "requirements.yaml"

		// Check if file already exists
		if _, err := os.Stat(filename); err == nil {
			os.Exit(0)
		}

		// Get current directory name
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error: failed to get current directory: %v\n", err)
			os.Exit(1)
		}

		// Create new project
		project := &types.Project{
			Name:         filepath.Base(currentDir),
			Requirements: []types.Requirement{},
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

		fmt.Println("Created requirements.yaml")

		// Check for OpenAI API key and inform user
		if os.Getenv("OPENAI_API_KEY") == "" {
			fmt.Println("\nFor requirement validation, set your OpenAI API key:")
			fmt.Println("  export OPENAI_API_KEY=\"your-api-key-here\"")
			fmt.Println("\n--no-validate will be used when adding requirements")
		}
	},
}
