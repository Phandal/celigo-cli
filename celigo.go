package main

import (
	"fmt"
	"github.com/Phandal/celigo-cli/arg"
	"os"
)

func usage() {
	fmt.Println("Usage:")
	fmt.Println("\tceligo <resource> <action> [options]")
}

func main() {
	cmd, err := arg.Parse(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		usage()
		os.Exit(1)
	}

	fmt.Printf("Resource: %v\n", cmd.Resource)
	fmt.Printf("Action: %v\n", cmd.Action)
}
