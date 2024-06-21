package celigo

import (
	"fmt"
)

type HelpAction struct {
	BaseAction
	resources *map[string]MappedResource
}

func (h *HelpAction) Execute() error {
	if isHelp, err := h.Parse("Usage: celigo-cli help\n"); err != nil {
		return err
	} else if isHelp {
		return nil
	}

	fmt.Printf("Usage: celigo-cli <resource> <action> [options]\n\n")
	PrintResources(*h.resources)
	return nil
}

func newHelpAction(args []string, resources *map[string]MappedResource) ActionExecuter {
	action := &HelpAction{
		BaseAction: BaseAction{
			usage: "show this help message",
			args:  args,
		},
		resources: resources,
	}

	return action
}

func NewHelpResource(usage string, cmd *Command) Resource {
	res := Resource{
		usage:   usage,
		actions: make(map[string]ActionExecuter),
	}

	res.newAction("help", newHelpAction(cmd.args, &cmd.mappedResources))

	return res
}
