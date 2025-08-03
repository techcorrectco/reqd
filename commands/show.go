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
			requirement := findRequirement(project.Requirements, requirementID)
			if requirement == nil {
				fmt.Fprintf(os.Stderr, "Error: Requirement '%s' not found\n", requirementID)
				os.Exit(1)
			}
			showRequirement(requirement, 0)
		} else {
			// Show entire list of requirements
			showRequirements(project.Requirements, 0)
		}
	},
}

// showRequirements renders a list of requirements with 2-space indentation
func showRequirements(requirements []types.Requirement, indentLevel int) {
	for _, req := range requirements {
		showRequirement(&req, indentLevel)
	}
}

// showRequirement renders a single requirement and its children with 2-space indentation
func showRequirement(req *types.Requirement, indentLevel int) {
	indent := ""
	for i := 0; i < indentLevel; i++ {
		indent += "  "
	}

	fmt.Printf("%s%s: %s\n", indent, req.ID, req.Title)

	if len(req.Children) > 0 {
		showRequirements(req.Children, indentLevel+1)
	}
}
