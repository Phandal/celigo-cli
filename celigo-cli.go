package main

import (
	"fmt"
	"os"

	"github.com/Phandal/celigo-cli/celigo"
)

func printError(name string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
}

func main() {
	programName := os.Args[0]

	cmd, err := celigo.NewCommand(os.Args)
	if err != nil {
		printError(programName, err)
		os.Exit(1)
	}

	cmd.NewResource("help", "shows this help message", celigo.NewHelpResource)
	// cmd.NewResource("script", "list, create, fetch, update, remove script", NewScriptResource)

	if err := cmd.Execute(); err != nil {
		printError(programName, err)
		os.Exit(1)
	}
}
