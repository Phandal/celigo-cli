package arg

import (
	"errors"
	"fmt"
)

type Command struct {
	Resource string
	Action   string
	Args     []string
}

type Value interface {
	String() string
	Set(string) error
}

type intValue int
type boolValue bool

// NOTE: trying to make this generic to handle all types
func newIntValue[T *int | *bool](v T, d T) *Value {
	switch v := any(v).(type) {
	case *int:
		*v = (*int)(d)
		return (*intValue)(v)
	case *bool:
		*v = d

	}
}

func (c *Command) RegisterInt(v *int, short string, long string, name string, usage string, defaultValue int) {
	c.NewFlag(newIntValue(v, defaultValue), short, long, name, usage)
}

func (c *Command) RegisterBool(v *bool, short string, long string, name string, usage string, defaultValue bool) {
	c.NewFlag(newBoolValue(v, defaultValue), short, long, name, usage)
}

func NewCommand(args []string) (Command, error) {
	if len(args) == 2 {
		if args[1] == "-h" || args[1] == "--help" || args[1] == "help" {
			return Command{Resource: "help", Action: ""}, nil
		}
	}

	if len(args) < 3 {
		return Command{}, errors.New("Invalid Number of Arguments")
	}

	cmd := Command{
		Resource: args[1],
		Action:   args[2],
		Args:     args[3:],
	}

	return cmd, nil
}
