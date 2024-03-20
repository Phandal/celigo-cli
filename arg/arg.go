package arg

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

type Command struct {
	Resource string
	Action   string
	Options  map[string]string
}

func (cmd *Command) RegisterFlag(o *string, short string, long string, name string, usage string, defValue string) error {
	valueShort, presentShort := cmd.Options[short]
	valueLong, presentLong := cmd.Options[long]

	if presentShort {
		*o = valueShort
	} else if presentLong {
		*o = valueLong
	} else if defValue != "" {
		*o = defValue
	} else {
		return fmt.Errorf("Missing argument: %s", name)
	}

	return nil
}

func parseFlags(flags []string) map[string]string {
	var parsedFlags = make(map[string]string)
	var currentArg string
	for _, v := range flags {
		if currentArg == "" {
			currentArg = v
			continue
		} else if strings.HasPrefix(currentArg, "-") || strings.HasPrefix(currentArg, "--") {
			parsedFlags[currentArg] = true
		} else {
			parsedFlags[currentArg] = v
			currentArg = ""
		}
	}

	return parsedFlags
}

func Parse(args []string) (Command, error) {
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
		Options:  parseFlags(args[3:]),
	}

	os.Args = os.Args[2:]

	return cmd, nil
}
