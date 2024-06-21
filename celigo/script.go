package celigo

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
	List   int = 200
	Create int = 201
	Fetch  int = 200
	Update int = 200
	Remove int = 204
)

const relativeUrl = "/scripts"
const filenameSeperator = "__"

type Script struct {
	Id                                     string `json:"_id"`
	LastModified                           string `json:"lastModified,omitempty"`
	CreatedAt                              string `json:"createdAt,omitempty"`
	Name                                   string `json:"name"`
	Description                            string `json:"description,omitempty"`
	Sandbox                                bool   `json:"sandbox,omitempty"`
	PostResponseHookToProcessOnChildRecord bool   `json:"postResponseHookToProcessOnChildRecord,omitempty"`
	Content                                string `json:"content,omitempty"`
}

func parseFilename(filename string) (string, string, error) {
	parsedFile := strings.Split(filename, filenameSeperator)

	if len(parsedFile) != 2 {
		return "", "", fmt.Errorf("Invalid Script Name: %s", filename)
	}

	return parsedFile[0], strings.TrimSuffix(parsedFile[1], filepath.Ext(parsedFile[1])), nil
}

func readScriptFile(scriptName string) ([]byte, error) {
	var file *os.File
	var content []byte
	var err error

	if file, err = os.Open(scriptName); err != nil {
		return []byte{}, err
	}

	if content, err = io.ReadAll(file); err != nil {
		return []byte{}, err
	}

	return content, nil
}

type ScriptHelpAction BaseHelpAction

func (s *ScriptHelpAction) Execute() error {
	if isHelp, err := s.Parse("Usage: celigo-cli script help\n"); err != nil {
		return err
	} else if isHelp {
		return nil
	}

	fmt.Printf("Script Resource Usage: celigo-cli script <action> [options]\n\n")
	PrintActions(*s.actions)
	return nil
}

func newScriptHelpAction(args []string, actions *map[string]ActionExecuter) *ScriptHelpAction {
	action := ScriptHelpAction{
		BaseAction: BaseAction{
			usage: "show this help message",
			args:  args,
		},
		actions: actions,
	}

	return &action
}

type ScriptListAction struct {
	BaseAction
}

func newScriptListAction(args []string) *ScriptListAction {
	action := ScriptListAction{
		BaseAction: BaseAction{
			usage: "list all scripts in Celigo",
			args:  args,
		},
	}

	return &action
}

func (s *ScriptListAction) Execute() error {
	var scripts []Script

	if isHelp, err := s.Parse("Usage: celigo-cli script list\n"); err != nil {
		return err
	} else if isHelp {
		return nil
	}

	if err := ExecuteGet(relativeUrl, List, &scripts); err != nil {
		return fmt.Errorf("failed to list scripts: %w", err)
	}

	fmt.Printf("%-32s%s\n", "ID", "NAME")
	for _, v := range scripts {
		fmt.Printf("%s\t%s\n", v.Id, v.Name)
	}

	return nil
}

type ScriptCreateAction struct {
	BaseAction
	title      string
	outputPath string
}

func newScriptCreateAction(args []string) *ScriptCreateAction {
	action := ScriptCreateAction{
		BaseAction: BaseAction{
			usage: "create a script in Celigo",
			args:  args,
			flags: make(map[string]*Flag, 5),
		},
	}

	action.RegisterString(&action.title, "t", "title", "title of the script to create", "", true)
	action.RegisterString(&action.outputPath, "o", "output", "output path to write the script file", "", false)

	return &action
}

func (s *ScriptCreateAction) Execute() error {
	if isHelp, err := s.Parse("Usage: celigo-cli script create [options]\n\n"); err != nil {
		return err
	} else if isHelp {
		return nil
	}

	script := Script{
		Name: s.title,
	}

	if err := ExecutePost(relativeUrl, &script, Create, &script); err != nil {
		return fmt.Errorf("Failed to create script: %s", err)
	}

	fmt.Printf("Successfully Created Script:\n%s\t %s\n", script.Id, script.Name)

	// TODO: Move to function that checks for overwrites and --force flag
	if s.outputPath != "" {
		filename := script.Name + filenameSeperator + script.Id + ".js"
		if _, err := os.Stat(s.outputPath); err != nil {
			return err
		}

		filepath := path.Join(s.outputPath, filename)
		err := os.WriteFile(filepath, []byte(script.Content), 0660)
		if err != nil {
			return err
		}
		fmt.Printf("Wrote Contents to file:\n%s\n", filepath)
	}

	return nil
}

type ScriptFetchAction struct {
	BaseAction
	id         string
	force      bool
	outputPath string
}

func newScriptFetchAction(args []string) *ScriptFetchAction {
	action := ScriptFetchAction{
		BaseAction: BaseAction{
			usage: "fetch a script from Celigo",
			args:  args,
			flags: make(map[string]*Flag, 5),
		},
	}

	action.RegisterString(&action.id, "i", "id", "the ID of the script to fetch", "", true)
	action.RegisterBool(&action.force, "f", "force", "overwrites the local script file", false, false)
	action.RegisterString(&action.outputPath, "o", "output", "output path to write the script file", "", false)

	return &action
}

func (s *ScriptFetchAction) Execute() error {
	var script Script

	if isHelp, err := s.Parse("Usage: celigo-cli script fetch [options]\n\n"); err != nil {
		return err
	} else if isHelp {
		return nil
	}

	if err := ExecuteGet(relativeUrl+"/"+s.id, Fetch, &script); err != nil {
		return fmt.Errorf("Failed to fetch script: %s", err)
	}

	// TODO: Move to function that checks for overwrites and --force flag
	if s.outputPath != "" {
		var filename = script.Name + "__" + script.Id + ".js"
		if _, err := os.Stat(s.outputPath); err != nil {
			return err
		}

		filepath := path.Join(s.outputPath, filename)
		err := os.WriteFile(filepath, []byte(script.Content), 0660)
		if err != nil {
			return err
		}
		fmt.Printf("Wrote Contents to file:\n%s\n", filepath)
	} else {
		fmt.Println(script.Content)
	}

	return nil
}

type ScriptUpdateAction struct {
	BaseAction
	filename string
	preview  bool
}

func newScriptUpdateAction(args []string) *ScriptUpdateAction {
	action := ScriptUpdateAction{
		BaseAction: BaseAction{
			usage: "update a script in Celigo",
			args:  args,
			flags: make(map[string]*Flag, 5),
		},
	}

	action.RegisterString(&action.filename, "i", "input", "path to script contents file", "", true)
	action.RegisterBool(&action.preview, "p", "preview", "preview the update without pushing changes to Celigo", false, false)

	return &action
}

func (s *ScriptUpdateAction) Execute() error {
	if isHelp, err := s.Parse("Usage: celigo-cli script update [options]\n\n"); err != nil {
		return err
	} else if isHelp {
		return nil
	}

	content, err := readScriptFile(s.filename)
	if err != nil {
		return err
	}

	scriptName, id, err := parseFilename(path.Base(s.filename))
	if err != nil {
		return err
	}

	if s.preview {
		fmt.Printf("Preview script update: %s\t%s\n", id, scriptName)
		return nil
	}

	script := Script{
		Id:      id,
		Name:    scriptName,
		Content: string(content),
	}

	if err = ExecutePut(relativeUrl+"/"+id, &script, Update, &script); err != nil {
		return fmt.Errorf("Failed to update script: %s", err)
	}

	fmt.Printf("Successfully Updated Script: %s\t%s\n", script.Id, script.Name)

	return nil
}

type ScriptRemoveAction struct {
	BaseAction
	id string
}

func newScriptRemoveAction(args []string) *ScriptRemoveAction {
	action := ScriptRemoveAction{
		BaseAction: BaseAction{
			usage: "remove a script from Celigo",
			args:  args,
			flags: make(map[string]*Flag, 5),
		},
	}

	action.RegisterString(&action.id, "i", "id", "id of the script to remove", "", true)

	return &action
}

func (s *ScriptRemoveAction) Execute() error {
	if isHelp, err := s.Parse("Usage: celigo-cli script update [options]\n\n"); err != nil {
		return err
	} else if isHelp {
		return nil
	}

	if err := ExecuteDelete(relativeUrl+"/"+s.id, Remove); err != nil {
		return fmt.Errorf("Failed to remove script: %s", err)
	}

	fmt.Printf("Successfully Removed Script with Id: %s", s.id)

	return nil
}

func NewScriptResource(usage string, cmd *Command) Resource {
	res := Resource{
		usage:   usage,
		actions: make(map[string]ActionExecuter),
	}

	res.newAction("help", newScriptHelpAction(cmd.args, &res.actions))
	res.newAction("list", newScriptListAction(cmd.args))
	res.newAction("create", newScriptCreateAction(cmd.args))
	res.newAction("fetch", newScriptFetchAction(cmd.args))
	res.newAction("update", newScriptUpdateAction(cmd.args))
	res.newAction("remove", newScriptRemoveAction(cmd.args))

	return res
}
