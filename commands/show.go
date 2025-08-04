package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/techcorrectco/reqd/internal/types"
)

var ShowCmd = &cobra.Command{
	Use:     "show [requirement_id]",
	Aliases: []string{"s"},
	Short:   "Display requirements",
	Long:    `Display project requirements or a specific requirement with its children.`,
	Args:    cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Load existing project
		project, err := types.LoadProject()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: No requirements.yaml found. Run 'reqd init' first.\n")
			os.Exit(1)
		}

		if len(args) > 0 {
			// Show specific requirement and its children
			requirementID := args[0]
			requirement := project.FindRequirement(requirementID)
			if requirement == nil {
				fmt.Fprintf(os.Stderr, "Error: Requirement '%s' not found\n", requirementID)
				os.Exit(1)
			}
			showRequirement(requirement)
		} else {
			// Show entire list of requirements
			showRequirements(project.Requirements)
		}
	},
}

// showRequirements renders a list of requirements
func showRequirements(requirements []types.Requirement) {
	for _, req := range requirements {
		showRequirement(&req)
	}
}

// showRequirement renders a single requirement and its children
func showRequirement(req *types.Requirement) {
	fmt.Printf("%s: %s\n", req.ID, req.Text)

	if len(req.Children) > 0 {
		showRequirements(req.Children)
	}
}
