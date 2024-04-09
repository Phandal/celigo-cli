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

const relativeUrl = "/scripts"
const filenameSeperator = "__"

func ExecuteScriptResource(cmd *Command) error {
	switch cmd.Action {
	case "list":
		return list(cmd)
	case "create":
		return create(cmd)
	case "fetch":
		return fetch(cmd)
	case "update":
		return update(cmd)
	case "remove":
		return remove(cmd)
	default:
		return fmt.Errorf("Unknown Action \"%s\"", cmd.Action)
	}
}

func list(_ *Command) error {
	var scripts []Script

	if err := ExecuteGet(relativeUrl, List, &scripts); err != nil {
		return fmt.Errorf("Failed to list scripts: %s", err)
	}

	fmt.Printf("%-32s%s\n", "ID", "NAME")
	for _, v := range scripts {
		fmt.Printf("%s\t%s\n", v.Id, v.Name)
	}

	return nil
}

func create(cmd *Command) error {
	var err error

	var title string
	cmd.RegisterString(&title, "t", "title", "title of the script to create", "", true)
	cmd.Parse()

	script := Script{
		Name: title,
	}

	if err = ExecutePost(relativeUrl, &script, Create, &script); err != nil {
		return fmt.Errorf("Failed to create script: %s", err)
	}

	fmt.Printf("Successfully Created Script:\n%s\t %s\n", script.Id, script.Name)

	return nil
}

func fetch(cmd *Command) error {
	var err error
	var script Script

	var id string
	var force bool
	var outputPath string

	cmd.RegisterString(&id, "i", "id", "the ID of the script to fetch", "", true)
	cmd.RegisterBool(&force, "f", "force", "overwrites the local script file", false, false)
	cmd.RegisterString(&outputPath, "o", "output", "output path to write the script file", "", false)

	err = cmd.Parse()
	if err != nil {
		return err
	}

	if err := ExecuteGet(relativeUrl+"/"+id, Fetch, &script); err != nil {
		return fmt.Errorf("Failed to fetch script: %s", err)
	}

	// TODO: Move to function that always writes to file, however checks for overwrite and only overwrites if the -o flag is present
	if outputPath != "" {
		var filename = script.Name + "__" + script.Id + ".js"
		if _, err := os.Stat(outputPath); err != nil {
			return err
		}

		filepath := path.Join(outputPath, filename)
		err = os.WriteFile(filepath, []byte(script.Content), 0660)
		if err != nil {
			return err
		}
		fmt.Printf("Wrote Contents to file:\n%s\n", filepath)
	} else {
		fmt.Println(script.Content)
	}

	return nil
}

func update(cmd *Command) error {
	var err error
	var scriptName string
	var id string
	var content []byte

	var filename string
	cmd.RegisterString(&filename, "i", "input", "path to script contents file", "", true)
	cmd.Parse()

	if content, err = readScriptFile(filename); err != nil {
		return err
	}

	if scriptName, id, err = parseFilename(path.Base(filename)); err != nil {
		return err
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

func remove(cmd *Command) error {
	var id string
	cmd.RegisterString(&id, "i", "id", "id of the script to remove", "", true)
	cmd.Parse()

	if err := ExecuteDelete(relativeUrl+"/"+id, Remove); err != nil {
		return fmt.Errorf("Failed to remove script: %s", err)
	}

	fmt.Printf("Successfully Removed Script with Id: %s", id)

	return nil
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
