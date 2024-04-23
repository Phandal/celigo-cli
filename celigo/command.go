package celigo

import (
	"errors"
	"fmt"
)

type Action interface {
	Execute() error
}

type Command struct {
	mappedResources map[string]MappedResource
	ResourceArg     string
	ActionArg       string
	Args            []string
}

type ResourceBuilder = func(*Command) Resource

type Resource struct {
	usage   string
	actions map[string]Action
}

type MappedResource struct {
	usage     string
	construct ResourceBuilder
}

func (c *Command) NewResource(name string, usage string, builder ResourceBuilder) {
	c.mappedResources[name] = MappedResource{usage: usage, construct: builder}
}

var missingActionErr = errors.New("missing action")

func invalidResourceErr(resourceName string) error {
	return fmt.Errorf("invalid resource name \"%s\"", resourceName)
}

func invalidActionErr(actionName string) error {
	return fmt.Errorf("invalid action name \"%s\"", actionName)
}

// Constructs a new Command struct
func NewCommand(args []string) (Command, error) {
	resources := make(map[string]MappedResource, 5)

	if len(args) < 2 {
		return Command{mappedResources: resources, ResourceArg: "help", ActionArg: "help"}, nil
	}
	if len(args) < 3 && (args[1] == "help" || args[1] == "--help" || args[1] == "-h") {
		return Command{mappedResources: resources, ResourceArg: "help", ActionArg: "help"}, nil
	}

	if len(args) < 3 {
		return Command{ResourceArg: args[1]}, missingActionErr
	}

	cmd := Command{
		mappedResources: make(map[string]MappedResource, 5),
		ResourceArg:     args[1],
		ActionArg:       args[2],
		Args:            args[3:],
	}

	return cmd, nil
}

func (c Command) Execute() error {
	mappedResource, exists := c.mappedResources[c.ResourceArg]
	if !exists {
		return invalidResourceErr(c.ResourceArg)
	}

	res := mappedResource.construct(&c)

	action, exists := res.actions[c.ActionArg]
	if !exists {
		return invalidActionErr(c.ActionArg)
	}

	return action.Execute()
}
