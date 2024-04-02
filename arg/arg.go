package arg

import (
	"errors"
	"fmt"
	"strconv"
)

type Command struct {
	Resource string
	Action   string
	Args     []string
	Flags    map[string]*Flag
}

type Flag struct {
	Value Value
	Short string
	Long  string
	Usage string
	Seen  bool
}

type Value interface {
	String() string
	Set(string) error
}

func numError(err error) error {
	ne, ok := err.(*strconv.NumError)
	if !ok {
		return err
	}
	if ne.Err == strconv.ErrSyntax {
		return errParse
	}
	if ne.Err == strconv.ErrRange {
		return errRange
	}

	return err
}

func missingError(name string) error {
	return fmt.Errorf("missing value for %s", name)
}

func unknownError(name string) error {
	return fmt.Errorf("unknown argument %s", name)
}

var errParse = errors.New("parse error")
var errRange = errors.New("value out of range")

// Int Flags
type intValue int

func newIntValue(v *int, d int) *intValue {
	*v = d
	return (*intValue)(v)
}

func (v *intValue) String() string {
	return strconv.FormatInt(int64(*v), 0)
}

func (v *intValue) Set(s string) error {
	n, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		return errParse
	}

	*v = intValue(n)
	return err
}

// Bool Flags
type boolValue bool

func newBoolValue(v *bool, d bool) *boolValue {
	*v = d
	return (*boolValue)(v)
}

func (v *boolValue) String() string {
	return strconv.FormatBool(bool(*v))
}

func (v *boolValue) Set(s string) error {
	n, err := strconv.ParseBool(s)
	if err != nil {
		return errParse
	}

	*v = boolValue(n)
	return err
}

// String Flags
type stringValue string

func newStringValue(v *string, d string) *stringValue {
	*v = d
	return (*stringValue)(v)
}

func (v *stringValue) String() string {
	return string(*v)
}

func (v *stringValue) Set(s string) error {
	*v = stringValue(s)
	return nil
}

// Command Functions
func (c *Command) newFlag(v Value, short string, long string, usage string, mandatory bool) {
	c.Flags[long] = &Flag{v, short, long, usage, !mandatory}
}

func (c *Command) RegisterInt(v *int, short string, long string, usage string, defaultValue int, mandatory bool) {
	c.newFlag(newIntValue(v, defaultValue), short, long, usage, mandatory)
}

func (c *Command) RegisterBool(v *bool, short string, long string, usage string, defaultValue bool, mandatory bool) {
	c.newFlag(newBoolValue(v, defaultValue), short, long, usage, mandatory)
}

func (c *Command) RegisterString(v *string, short string, long string, usage string, defaultValue string, mandatory bool) {
	c.newFlag(newStringValue(v, defaultValue), short, long, usage, mandatory)
}

func NewCommand(args []string) (Command, error) {
	if len(args) == 2 {
		if args[1] == "-h" || args[1] == "--help" || args[1] == "help" {
			return Command{Resource: "help"}, nil
		}
	}

	if len(args) < 3 {
		return Command{}, errors.New("not enough arguments")
	}

	cmd := Command{
		Resource: args[1],
		Action:   args[2],
		Args:     args[3:],
		Flags:    make(map[string]*Flag),
	}

	return cmd, nil
}

func (c *Command) parseOne() (bool, error) {
	// Make sure there are args to parse
	if len(c.Args) == 0 {
		return false, nil
	}

	// Needs to be at least 2 characters to be an arg
	s := c.Args[0]
	if len(s) < 2 || s[0] != '-' {
		return false, nil
	}

	count := 1
	if s[1] == '-' {
		count++
	}

	// Find matching flag in registered flags
	arg := s[count:]
	var flag *Flag
	var name string
	var found bool = false
	for name, flag = range c.Flags {
		if arg == flag.Short || arg == flag.Long {
			found = true
			break
		}
	}

	if !found {
		return false, unknownError(arg)
	}

	c.Args = c.Args[1:]

	// Handle bool flags
	if _, ok := flag.Value.(*boolValue); ok {
		flag.Seen = true
		flag.Value.Set("true")
		return true, nil
	}

	// Handle all other flags
	if len(c.Args) == 0 {
		return false, missingError(name)
	}

	value := c.Args[0]
	if value[0] == '-' {
		return false, missingError(name)
	}

	flag.Seen = true
	flag.Value.Set(c.Args[0])
	c.Args = c.Args[1:]
	return true, nil

}

func (c *Command) Parse() error {
	for {
		seen, err := c.parseOne()
		if err != nil {
			return err
		}

		if seen {
			continue
		}

		break
	}

	for k, v := range c.Flags {
		if !v.Seen {
			return fmt.Errorf("missing %s\n", k)
		}
	}

	return nil
}
