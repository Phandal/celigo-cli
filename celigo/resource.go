package celigo

import "fmt"

type ActionExecuter interface {
	Execute() error
	Usage() string
}

type BaseAction struct {
	usage string
	args  []string
	flags map[string]*Flag
}

func (b BaseAction) Usage() string {
	return b.usage
}

type BaseHelpAction struct {
	BaseAction
	actions *map[string]ActionExecuter
}

type Resource struct {
	usage   string
	actions map[string]ActionExecuter
	args    []string
}

func formatActionForHelpMessage(name string, action ActionExecuter) string {
	return fmt.Sprintf("  %-15s%s\n", name, action.Usage())
}

func (r *Resource) newAction(name string, action ActionExecuter) {
	r.actions[name] = action
}
