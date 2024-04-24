package celigo

import (
	"errors"
	"fmt"
	"os"
)

type ProjectHelpAction BaseHelpAction

func newProjectHelpAction(args []string, actions *map[string]ActionExecuter) *ProjectHelpAction {
	action := &ProjectHelpAction{
		BaseAction: BaseAction{
			usage: "show this help message",
			args:  args,
		},
		actions: actions,
	}

	return action
}

func (p *ProjectHelpAction) Execute() error {
	if isHelp, err := p.Parse("Usage: celigo-cli project help\n"); err != nil {
		return err
	} else if isHelp {
		return nil
	}

	fmt.Printf("Project Resource Usage: celigo-cli project <action> [options]\n")
	for name, action := range *p.actions {
		fmt.Printf("  %-15s%s\n", name, action.Usage())
	}
	return nil
}

type ProjectInitAction struct {
	BaseAction
	apiKey string
}

func (p *ProjectInitAction) Execute() error {
	if isHelp, err := p.Parse("Usage: celigo-cli project init [options]\n\n"); err != nil {
		return err
	} else if isHelp {
		return nil
	}

	// TODO move this to function to check if file exists
	if _, err := os.Stat(".env"); !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf(".env file already exists")
	}

	// TODO move to a function
	if p.apiKey == "" {
		fmt.Print("Please enter your api key: ")
		if _, err := fmt.Scanln(&p.apiKey); err != nil {
			return fmt.Errorf("failed to take user input: %w", err)
		}
	}

	err := os.WriteFile(".env", []byte(p.apiKey), 0660)
	if err != nil {
		return err
	}
	fmt.Printf("Created .env file\n")

	if _, err := os.Stat(".celigo-cli"); !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("project already initilized")
	}

	err = os.WriteFile(".celigo-cli", nil, 0660)
	if err != nil {
		return err
	}
	fmt.Printf("Project Initialized!\n")

	return nil
}

func newProjectInitAction(args []string) *ProjectInitAction {
	action := &ProjectInitAction{
		BaseAction: BaseAction{
			usage: "initialize a new project",
			args:  args,
			flags: make(map[string]*Flag, 5),
		},
	}

	action.RegisterString(&action.apiKey, "k", "key", "api key from celigo", "", false)

	return action
}

func NewProjectResource(usage string, cmd *Command) Resource {
	res := Resource{
		usage:   usage,
		actions: make(map[string]ActionExecuter),
	}

	res.newAction("help", newProjectHelpAction(cmd.args, &res.actions))
	res.newAction("init", newProjectInitAction(cmd.args))

	return res
}
