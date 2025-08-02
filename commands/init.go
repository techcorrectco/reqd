package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/techcorrectco/reqd/internal/types"
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
		var idPrefix = generateIDPrefix(defaultName)

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
		project := &types.Project{
			Name:         projectName,
			IDPrefix:     idPrefix,
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

		fmt.Printf("'%s' is ready for requirements\n", project.Name)
	},
}

// generateIDPrefix creates an ID prefix from a project name using acronym generation
func generateIDPrefix(projectName string) string {
	// Clean and normalize the input
	name := strings.TrimSpace(projectName)
	if name == "" {
		return ""
	}

	// Split on common delimiters
	words := strings.FieldsFunc(name, func(r rune) bool {
		return r == ' ' || r == '-' || r == '_' || r == '.'
	})

	// If only one word or no words, use fallback approach
	if len(words) <= 1 {
		maxLen := 4
		if len(name) < maxLen {
			maxLen = len(name)
		}
		return strings.ToUpper(name[:maxLen])
	}

	// Multiple words: generate acronym
	var prefix strings.Builder

	for _, word := range words {
		word = strings.TrimSpace(word)
		if len(word) > 0 {
			// Take first character of each word
			firstChar := rune(word[0])
			if unicode.IsLetter(firstChar) {
				prefix.WriteRune(unicode.ToUpper(firstChar))
			}
		}
	}

	result := prefix.String()

	// Fallback: if no letters found in acronym, use first 3-4 chars uppercased
	if len(result) == 0 {
		maxLen := 4
		if len(name) < maxLen {
			maxLen = len(name)
		}
		result = strings.ToUpper(name[:maxLen])
	}

	return result
}

