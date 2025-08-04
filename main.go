package main

import (
	"fmt"
	"os"

	"github.com/techcorrectco/reqd/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
