package celigo

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

func (r *Resource) newAction(name string, action ActionExecuter) {
	r.actions[name] = action
}
