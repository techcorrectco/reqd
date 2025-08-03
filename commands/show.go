package commands

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh"
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
	for _, req := range requirements {
		title := req.Title
		if len(req.Children) > 0 {
			title = "+ " + title
		}
		display := fmt.Sprintf("%s: %s", req.ID, title)
		options = append(options, huh.NewOption(display, req.ID))
	}

	// Add quit option
	options = append(options, huh.NewOption("Quit", "q"))

	var selected string
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

	// Handle quit
	if selected == "q" {
		return
	}

	// Find selected requirement
	var selectedReq *types.Requirement
	for i := range requirements {
		if requirements[i].ID == selected {
			selectedReq = &requirements[i]
			break
		}
	}

	if selectedReq != nil && len(selectedReq.Children) > 0 {
		// Re-render with children if they exist
		showRequirementsList(project, selectedReq)
	}
	// Do nothing if no children
}
