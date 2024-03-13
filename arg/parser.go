package arg

import (
	"errors"
)

type Flag struct {
	Name  string
	Value string
}

type Command struct {
	Resource string
	Action   string
	Options  []Flag
}

func parseFlags(flags []string) []Flag {
	var parsedFlags []Flag
	var currentArg string
	for _, v := range flags {
		if currentArg == "" {
			currentArg = v
			continue
		} else {
			parsedFlags = append(parsedFlags, Flag{Name: currentArg, Value: v})
			currentArg = ""
		}
	}

	return parsedFlags
}

func Parse(args []string) (Command, error) {
	if len(args) == 2 {
		if args[1] == "-h" || args[1] == "--help" {
			return Command{Resource: "help", Action: "", Options: []Flag{}}, nil
		}
	}

	if len(args) < 3 {
		return Command{}, errors.New("Invalid Number of Arguments")
	}

	cmd := Command{
		Resource: args[1],
		Action:   args[2],
		Options:  parseFlags(args[3:]),
	}

	return cmd, nil
}
