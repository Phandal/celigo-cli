package main

import (
	"fmt"
	"os"

	"github.com/Phandal/celigo-cli/celigo"
	"github.com/Phandal/celigo-cli/dotenv"
)

func printError(name string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
}

func main() {
	programName := os.Args[0]
	// TODO: Allow the user to specify and environment file
	dotenv.Parse("")

	cmd, err := celigo.NewCommand(os.Args)
	if err != nil {
		printError(programName, err)
		os.Exit(1)
	}

	cmd.NewResource("project", "manages projects", celigo.NewProjectResource)
	cmd.NewResource("help", "shows this help message", celigo.NewHelpResource)
	cmd.NewResource("script", "list, create, fetch, update, remove script", celigo.NewScriptResource)

	if err := cmd.Execute(); err != nil {
		printError(programName, err)
		os.Exit(1)
	}
}
