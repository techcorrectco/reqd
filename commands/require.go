package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/techcorrectco/reqd/internal/openai"
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
		noValidate, _ := cmd.Flags().GetBool("no-validate")
		noParentProposal, _ := cmd.Flags().GetBool("no-parent-proposal")

		// Load existing project
		project, err := types.LoadProject()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: No requirements.yaml found. Run 'reqd init' first.\n")
			os.Exit(1)
		}

		var finalTitle string
		// Auto-skip validation if no API key is set and --no-validate wasn't explicitly used
		if noValidate || os.Getenv("OPENAI_API_KEY") == "" {
			// Skip validation, use original title
			finalTitle = requirementTitle
		} else {
			// Validate requirement with OpenAI
			finalTitle, err = validateRequirement(requirementTitle)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
				fmt.Fprintf(os.Stderr, "Proceeding with original requirement...\n")
				finalTitle = requirementTitle
			}
		}

		// If no parent ID provided and parent proposal not disabled, ask if user wants a parent proposed
		if parentID == "" && !noParentProposal && os.Getenv("OPENAI_API_KEY") != "" {
			proposedParent, err := proposeRequirementParent(finalTitle, project)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
			} else if proposedParent != "" {
				parentID = proposedParent
			}
		}

		// Generate new requirement
		newReq := createRequirement(finalTitle, parentID, project)

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

		fmt.Printf("\n%s: %s\n", newReq.ID, newReq.Text)
	},
}

func init() {
	RequireCmd.Flags().StringP("parent", "p", "", "Parent requirement ID")
	RequireCmd.Flags().BoolP("no-validate", "V", false, "Skip OpenAI validation of the requirement")
	RequireCmd.Flags().BoolP("no-parent-proposal", "P", false, "Skip proposing a parent for this requirement")
}

// createRequirement generates a new requirement with proper ID
func createRequirement(title, parentID string, project *types.Project) types.Requirement {
	var id string

	if parentID == "" {
		// Top-level requirement: use sequence number only
		id = fmt.Sprintf("%d", len(project.Requirements)+1)
	} else {
		// Child requirement: find parent and generate child ID
		parent := project.FindRequirement(parentID)
		if parent != nil {
			id = fmt.Sprintf("%s.%d", parentID, len(parent.Children)+1)
		} else {
			// Fallback if parent not found
			id = fmt.Sprintf("%d", len(project.Requirements)+1)
		}
	}

	return types.Requirement{
		ID:       id,
		Text:     title,
		Children: []types.Requirement{},
	}
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

// validateRequirement validates a requirement using OpenAI and returns the final title to use
func validateRequirement(input string) (string, error) {
	fmt.Println("Reviewing...")

	validation, err := openai.ValidateRequirement(input)
	if err != nil {
		return "", err
	}

	// Display validation results
	fmt.Printf("\nInput:\n%s\n\n", input)

	if len(validation.Problems) > 0 {
		fmt.Println("Issues:")
		for _, problem := range validation.Problems {
			fmt.Printf("- %s\n", problem)
		}
		fmt.Println()
	}

	fmt.Printf("Recommended:\n%s\n\n", validation.Recommended)

	// Ask user if they want to accept recommended changes
	fmt.Print("Accept recommended changes? [Y/n]: ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read user input: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))

	// Default to "yes" if empty response or "y"
	if response == "" || response == "y" || response == "yes" {
		return validation.Recommended, nil
	}

	return input, nil
}

// proposeRequirementParent asks user if they want a parent proposed and handles the proposal
func proposeRequirementParent(requirement string, project *types.Project) (string, error) {
	// Ask user if they want a parent proposed (default to yes)
	fmt.Print("Would you like a parent proposed for this requirement? [Y/n]: ")
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read user input: %w", err)
	}

	response = strings.TrimSpace(strings.ToLower(response))

	// Default to "yes" if empty response or "y"
	if response != "n" && response != "no" {
		// Auto-skip parent proposal if no API key is set
		if os.Getenv("OPENAI_API_KEY") == "" {
			fmt.Println("Skipping parent proposal (no OPENAI_API_KEY set)")
			return "", nil
		}

		// Get all branch requirements
		branches := project.GetBranches()
		if len(branches) == 0 {
			fmt.Println("No existing requirements with children found to use as parents.")
			return "", nil
		}

		// Get parent proposal from OpenAI
		proposal, err := openai.ProposeParent(requirement, branches)
		if err != nil {
			return "", fmt.Errorf("failed to get parent proposal: %w", err)
		}

		if proposal.ProposedParent != nil && *proposal.ProposedParent != "" {
			proposedParent := project.FindRequirement(*proposal.ProposedParent)
			if proposedParent != nil {
				fmt.Printf("\nSuggested parent: %s\n", proposedParent.DisplayFormat())
				fmt.Print("Accept suggested parent? [Y/n]: ")
				
				acceptResponse, err := reader.ReadString('\n')
				if err != nil {
					return "", fmt.Errorf("failed to read user input: %w", err)
				}

				acceptResponse = strings.TrimSpace(strings.ToLower(acceptResponse))
				
				// Default to "yes" if empty response or "y"
				if acceptResponse == "" || acceptResponse == "y" || acceptResponse == "yes" {
					return *proposal.ProposedParent, nil
				}
			}
		} else {
			fmt.Println("No suitable parent found.")
		}
	}

	return "", nil
}
