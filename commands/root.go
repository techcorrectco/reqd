package commands

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "reqd",
	Short: "Manage the complexity of your Product Requirements Document (PRD)",
	Long: `reqd is a CLI tool designed to help you manage the complexity 
of your Product Requirements Document (PRD). It provides commands 
to create, organize, and maintain your requirements effectively.`,
}

func init() {
	RootCmd.AddCommand(InitCmd)
	RootCmd.AddCommand(RequireCmd)
	RootCmd.AddCommand(ShowCmd)
}
