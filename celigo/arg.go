package celigo

import (
	"errors"
	"fmt"
	"strconv"
)

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

func (f *Flag) formatForHelpMessage(name string) string {
	return fmt.Sprintf("  -%-2s, --%-10s%s\n", f.Short, f.Long, f.Usage)
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
func (b *BaseAction) newFlag(v Value, short string, long string, usage string, mandatory bool) {
	b.flags[long] = &Flag{v, short, long, usage, !mandatory}
}

/* Registers an int arg to get parsed */
func (b *BaseAction) RegisterInt(v *int, short string, long string, usage string, defaultValue int, mandatory bool) {
	b.newFlag(newIntValue(v, defaultValue), short, long, usage, mandatory)
}

/* Registers a bool arg to get parsed */
func (b *BaseAction) RegisterBool(v *bool, short string, long string, usage string, defaultValue bool, mandatory bool) {
	b.newFlag(newBoolValue(v, defaultValue), short, long, usage, mandatory)
}

/* Registers a string arg to get parsed */
func (b *BaseAction) RegisterString(v *string, short string, long string, usage string, defaultValue string, mandatory bool) {
	b.newFlag(newStringValue(v, defaultValue), short, long, usage, mandatory)
}

func (b *BaseAction) parseOne() (bool, bool, error) {
	// Make sure there are args to parse
	if len(b.args) == 0 {
		return false, false, nil
	}
	s := b.args[0]

	// Check for Help
	if s == "help" || s == "--help" || s == "-h" {
		return false, true, nil
	}

	// Needs to be at least 2 characters to be an arg
	if len(s) < 2 || s[0] != '-' {
		return false, false, nil
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
	for name, flag = range b.flags {
		if arg == flag.Short || arg == flag.Long {
			found = true
			break
		}
	}

	if !found {
		return false, false, unknownError(arg)
	}

	b.args = b.args[1:]

	// Handle bool flags
	if _, ok := flag.Value.(*boolValue); ok {
		flag.Seen = true
		flag.Value.Set("true")
		return true, false, nil
	}

	// Handle all other flags
	if len(b.args) == 0 {
		return false, false, missingError(name)
	}

	value := b.args[0]
	if value[0] == '-' {
		return false, false, missingError(name)
	}

	flag.Seen = true
	flag.Value.Set(b.args[0])
	b.args = b.args[1:]
	return true, false, nil

}

func (b *BaseAction) Help(usageMessage string) {
	fmt.Printf("%s", usageMessage)
	PrintFlags(b.flags)
}

/* Parses registered Action flags into their respective variables */
func (b *BaseAction) Parse(usageMessage string) (bool, error) {
	for {
		seen, isHelp, err := b.parseOne()
		if err != nil {
			return false, err
		}

		if isHelp {
			b.Help(usageMessage)
			return true, nil
		}

		if seen {
			continue
		}

		break
	}

	if len(b.args) != 0 {
		return false, fmt.Errorf("unknown options: %v\n", b.args)
	}

	for k, v := range b.flags {
		if !v.Seen {
			return false, fmt.Errorf("missing option \"%s\"\n", k)
		}
	}

	return false, nil
}
