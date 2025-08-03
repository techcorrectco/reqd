package commands

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/techcorrectco/reqd/internal/types"
)

var ShowCmd = &cobra.Command{
	Use:     "show [requirement_id]",
	Aliases: []string{"s"},
	Short:   "Display requirements in an interactive list",
	Long:    `Display project requirements or children of a specific requirement using an interactive selection list.`,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Load existing project
		project, err := types.LoadProject()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: No requirements.yaml found. Run 'reqd init' first.\n")
			os.Exit(1)
		}

		var startingRequirement *types.Requirement
		if len(args) > 0 {
			requirementID := args[0]
			startingRequirement = findRequirement(project.Requirements, requirementID)
			if startingRequirement == nil {
				fmt.Fprintf(os.Stderr, "Error: Requirement '%s' not found\n", requirementID)
				os.Exit(1)
			}
		}

		showRequirementsList(project, startingRequirement)
	},
}

// showRequirementsList displays an interactive list and handles navigation
func showRequirementsList(project *types.Project, currentReq *types.Requirement) {
	showRequirementsListWithParent(project, currentReq, nil)
}

// showRequirementsListWithParent displays an interactive list with parent context for back navigation
func showRequirementsListWithParent(project *types.Project, currentReq *types.Requirement, parentReq *types.Requirement) {
	var title string
	var requirements []types.Requirement

	if currentReq == nil {
		// Show project-level requirements
		title = project.Name
		requirements = project.Requirements
	} else {
		// Show children of specified requirement
		title = fmt.Sprintf("%s: %s", currentReq.ID, currentReq.Title)
		requirements = currentReq.Children
	}

	if len(requirements) == 0 {
		fmt.Printf("No requirements found for %s\n", title)
		return
	}

	// Prepare options for huh select
	var options []huh.Option[string]

	// Add back navigation option if we're viewing children (not at root level)
	if currentReq != nil {
		options = append(options, huh.NewOption("..", ".."))
	}

	for _, req := range requirements {
		title := req.Title
		if len(req.Children) > 0 {
			title = "+ " + title
		}
		display := fmt.Sprintf("%s: %s", req.ID, title)
		options = append(options, huh.NewOption(display, req.ID))
	}

	var selected string

	// First, get the requirement selection
	err := huh.NewSelect[string]().
		Title(title).
		Options(options...).
		Value(&selected).
		WithTheme(huh.ThemeBase()).
		Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Find selected requirement
	var selectedReq *types.Requirement
	for i := range requirements {
		if requirements[i].ID == selected {
			selectedReq = &requirements[i]
			break
		}
	}

	// Handle back navigation
	if selected == ".." {
		if parentReq != nil {
			// Find the parent of parentReq to continue the chain
			grandParentReq := findParentRequirement(project.Requirements, parentReq.ID)
			showRequirementsListWithParent(project, parentReq, grandParentReq)
		} else {
			// Go back to root level
			showRequirementsListWithParent(project, nil, nil)
		}
		return
	}

	// If requirement has children, show them directly
	if len(selectedReq.Children) > 0 {
		showRequirementsListWithParent(project, selectedReq, currentReq)
		return
	}

	// If no children, show confirm for Edit/Quit
	shouldEdit := true
	err = huh.NewConfirm().
		WithButtonAlignment(lipgloss.Left).
		Title(fmt.Sprintf("%s: %s", selectedReq.ID, selectedReq.Title)).
		Affirmative("Edit").
		Negative("Quit").
		Value(&shouldEdit).
		WithTheme(huh.ThemeBase()).
		Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if shouldEdit {
		editRequirement(project, selectedReq)
	}
}

// editRequirement allows editing a requirement's title
func editRequirement(project *types.Project, req *types.Requirement) {
	newTitle := req.Title

	err := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title(fmt.Sprintf("Editing %s", req.ID)).
				Value(&newTitle).
				Placeholder(req.Title),
		),
	).WithTheme(huh.ThemeBase()).Run()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Update the requirement title if it changed
	if newTitle != req.Title && newTitle != "" {
		req.Title = newTitle

		// Save the updated project
		err = project.Save()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error saving project: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Updated requirement %s: %s\n", req.ID, req.Title)
	}
}

// findParentRequirement finds the parent requirement of a given requirement ID
func findParentRequirement(requirements []types.Requirement, targetID string) *types.Requirement {
	for i := range requirements {
		// Check if any of this requirement's children match the target
		for _, child := range requirements[i].Children {
			if child.ID == targetID {
				return &requirements[i]
			}
		}
		// Recursively search in children
		if parent := findParentRequirement(requirements[i].Children, targetID); parent != nil {
			return parent
		}
	}
	return nil
}
