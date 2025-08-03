package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/techcorrectco/reqd/internal/types"
)

var RequireCmd = &cobra.Command{
	Use:     "require [requirement title]",
	Aliases: []string{"r"},
	Short:   "Document a new system requirement",
	Long:    `Add a new requirement to the project. Generates an ID and adds it to the requirements hierarchy.`,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		requirementTitle := args[0]
		parentID, _ := cmd.Flags().GetString("parent")

		// Load existing project
		project, err := types.LoadProject()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: No requirements.yaml found. Run 'reqd init' first.\n")
			os.Exit(1)
		}

		// Generate new requirement
		newReq := createRequirement(requirementTitle, parentID, project)

		// Add requirement to project
		if parentID == "" {
			// Add as top-level requirement
			project.Requirements = append(project.Requirements, newReq)
		} else {
			// Find parent and add as child
			if !addChildRequirement(project.Requirements, parentID, newReq) {
				fmt.Fprintf(os.Stderr, "Error: Parent requirement '%s' not found\n", parentID)
				os.Exit(1)
			}
		}

		// Save project
		if err := project.Save(); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving requirements: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Added requirement %s: %s\n", newReq.ID, newReq.Title)
	},
}

func init() {
	RequireCmd.Flags().StringP("parent", "p", "", "Parent requirement ID")
}

// createRequirement generates a new requirement with proper ID
func createRequirement(title, parentID string, project *types.Project) types.Requirement {
	var id string

	if parentID == "" {
		// Top-level requirement: use sequence number only
		id = fmt.Sprintf("%d", len(project.Requirements)+1)
	} else {
		// Child requirement: find parent and generate child ID
		parent := findRequirement(project.Requirements, parentID)
		if parent != nil {
			id = fmt.Sprintf("%s.%d", parentID, len(parent.Children)+1)
		} else {
			// Fallback if parent not found
			id = fmt.Sprintf("%d", len(project.Requirements)+1)
		}
	}

	return types.Requirement{
		ID:          id,
		Title:       title,
		Keyword:     "", // Could be extracted from title or left empty
		Description: "", // Could be added later
		Children:    []types.Requirement{},
	}
}

// findRequirement recursively searches for a requirement by ID
func findRequirement(requirements []types.Requirement, id string) *types.Requirement {
	for i := range requirements {
		req := requirements[i]
		if req.ID == id {
			return &requirements[i]
		}
		if found := findRequirement(req.Children, id); found != nil {
			return found
		}
	}
	return nil
}

// addChildRequirement finds the parent and adds the child requirement
func addChildRequirement(requirements []types.Requirement, parentID string, child types.Requirement) bool {
	for i := range requirements {
		req := &requirements[i]
		if req.ID == parentID {
			req.Children = append(req.Children, child)
			return true
		}
		if addChildRequirement(req.Children, parentID, child) {
			return true
		}
	}
	return false
}