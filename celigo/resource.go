package celigo

import "fmt"

type HelpAction struct {
	usage     string
	resources map[string]MappedResource
}

func (h HelpAction) Execute() error {
	fmt.Printf("Usage: celigo-cli <resource> <action> [options]\n\n")
	for name, res := range h.resources {
		fmt.Printf("  %-15s%s\n", name, res.usage)
	}
	return nil
}

func NewHelpResource(cmd *Command) Resource {
	actions := map[string]Action{
		"help": HelpAction{usage: "show this help message", resources: cmd.mappedResources},
	}

	return Resource{
		usage:   "show this help message",
		actions: actions,
	}
}
