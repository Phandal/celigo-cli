package celigo

import (
	"fmt"
	"slices"
)

func PrintResources(resources map[string]MappedResource) {
	names := make([]string, 0)
	for name, _ := range resources {
		names = append(names, name)
	}
	slices.Sort(names)

	for _, name := range names {
		fmt.Printf(resources[name].formatForHelp(name))
	}
}

func PrintActions(actions map[string]ActionExecuter) {
	names := make([]string, 0)
	for name, _ := range actions {
		names = append(names, name)
	}
	slices.Sort(names)

	for _, name := range names {
		action := actions[name]
		fmt.Printf(formatActionForHelpMessage(name, action))
	}
}

func PrintFlags(flags map[string]*Flag) {
	names := make([]string, 0)
	for name, _ := range flags {
		names = append(names, name)
	}
	slices.Sort(names)

	for _, name := range names {
		fmt.Print(flags[name].formatForHelpMessage(name))
	}
}
