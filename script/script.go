package script

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/Phandal/celigo-cli/arg"
	"github.com/Phandal/celigo-cli/celigo"
)

const (
	List   int = 200
	Fetch  int = 200
	Update int = 200
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

func Execute(cmd *arg.Command) error {
	switch cmd.Action {
	case "list":
		return list(cmd)
	case "fetch":
		return fetch(cmd)
	case "update":
		return update(cmd)
	default:
		return fmt.Errorf("Unknown Action \"%s\"", cmd.Action)
	}
}

func list(_ *arg.Command) error {
	var scripts []Script

	if err := celigo.ExecuteGet(relativeUrl, List, &scripts); err != nil {
		return fmt.Errorf("Failed to list scripts: %s", err)
	}

	for _, v := range scripts {
		fmt.Printf("%s\t%s\n", v.Id, v.Name)
	}

	return nil
}

func fetch(cmd *arg.Command) error {
	var err error
	var script Script

	// var id = "65f10091892c590e57254963" // TODO: Find a way to search through the flags for this value
	// var shouldWrite = true              // TODO: make boolean flags work in the parser

	var id string
	var force bool

	cmd.RegisterFlag(&id, "-i", "--id", "id", "the ID of the script to fetch", "")
	cmd.RegisterFlag(&force, "-f", "--force", "force", "overwrites the local script file", false)

	cmd.ParseFlage()

	if err := celigo.ExecuteGet(relativeUrl+"/"+id, Fetch, &script); err != nil {
		return fmt.Errorf("Failed to fetch script: %s", err)
	}

	fmt.Printf("ID: %s\nName: %s\nDescription: %s\nLast Modified Date: %s\n", script.Id, script.Name, script.Description, script.LastModified)

	// TODO: Move to function that always writes to file, however checks for overwrite and only overwrites if the -o flag is present
	if force == true {
		var filename = script.Name + "__" + script.Id + ".js"
		err = os.WriteFile(filename, []byte(script.Content), 0660)
		if err != nil {
			return err
		}
		fmt.Printf("Wrote Contents to file: %s\n", filename)
	}

	return nil
}

func update(_ *arg.Command) error {
	var err error
	var scriptName string
	var id string
	var content []byte

	var filename = "Test Script 1__65f10091892c590e57254963.js" // TODO: Find a way to get this from the user

	if content, err = readScriptFile(filename); err != nil {
		return err
	}

	if scriptName, id, err = parseFilename(filename); err != nil {
		return err
	}

	script := Script{
		Id:      id,
		Name:    scriptName,
		Content: string(content),
	}

	if err = celigo.ExecutePut(relativeUrl+"/"+id, &script, Update, &script); err != nil {
		return fmt.Errorf("Failed to update script: %s", err)
	}

	fmt.Printf("Successfully Updated Script: %s\t%s\n", script.Id, script.Name)

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
