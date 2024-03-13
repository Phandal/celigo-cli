package main

import (
	"fmt"
	"os"

	"github.com/Phandal/celigo-cli/arg"
	"github.com/Phandal/celigo-cli/script"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Printf("\t%v <resource> <action> [options]\n", os.Args[0])
}

func run(cmd *arg.Command) error {
	switch cmd.Resource {
	case "help":
		usage()
		return nil
	case "script":
		return script.Execute(cmd)
	default:
		return fmt.Errorf("Unknown Resource Type \"%v\"\n", cmd.Resource)
	}
}

func main() {
	var cmd arg.Command
	var err error

	cmd, err = arg.Parse(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, ": %v\n", err)
		usage()
		os.Exit(1)
	}

	err = run(&cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v: %v\n", os.Args[0], err)
		os.Exit(1)
	}
}
