package arg

import (
	"errors"
)

type Flag struct {
	name  string
	value string
}

type Command struct {
	Resource string
	Action   string
	Options  []Flag
}

func Parse(args []string) (Command, error) {
	if len(args) < 3 {
		return Command{}, errors.New("Invalid Number of Arguements")
	}

	cmd := Command{
		Resource: args[1],
		Action:   args[2],
		Options:  []Flag{},
	}

	return cmd, nil
}
