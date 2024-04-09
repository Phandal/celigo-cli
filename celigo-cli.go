package main

import (
	"fmt"
	"os"

	"github.com/Phandal/celigo-cli/celigo"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Printf("\t%s <resource> <action> [options]\n", os.Args[0])
}

func run(cmd *celigo.Command) error {
	switch cmd.Resource {
	case "help":
		usage()
		return nil
	case "script":
		return celigo.ExecuteScriptResource(cmd)
	default:
		return fmt.Errorf("Unknown Resource Type \"%s\"", cmd.Resource)
	}
}

func main() {
	var cmd celigo.Command
	var err error

	cmd, err = celigo.NewCommand(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		usage()
		os.Exit(1)
	}

	err = run(&cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", os.Args[0], err)
		os.Exit(1)
	}
}
