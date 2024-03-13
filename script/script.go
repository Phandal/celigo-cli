package script

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/Phandal/celigo-cli/arg"
	"github.com/Phandal/celigo-cli/util"
)

type Script struct {
	Id                                     string `json:"_id"`
	LastModified                           string `json:"lastModified"`
	CreatedAt                              string `json:"createdAt"`
	Name                                   string `json:"name"`
	Description                            string `json:"description"`
	Sandbox                                bool   `json:"sandbox"`
	PostResponseHookToProcessOnChildRecord bool   `json:"postResponseHookToProcessOnChildRecord"`
	Content                                string `json:"content"`
}

const relativeUrl = "/scripts"

func Execute(cmd *arg.Command) error {
	switch cmd.Action {
	case "list":
		return list(cmd)
	case "fetch":
		return fetch(cmd)
	default:
		return fmt.Errorf("Unknown Action \"%v\"\n", cmd.Action)
	}
}

func list(_ *arg.Command) error {
	var client http.Client
	var err error
	var listSuccessStatusCode = 200
	var req *http.Request
	var res *http.Response
	var scripts []Script

	req, err = util.BuildRequest("GET", util.BaseUrl+relativeUrl, nil)
	client = http.Client{}

	res, err = client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != listSuccessStatusCode {
		return fmt.Errorf("List Status Code %v\n", res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(&scripts)

	if err != nil {
		return err
	}

	for _, v := range scripts {
		fmt.Printf("%v\t %v\n", v.Id, v.Name)
	}

	return nil
}

func fetch(_ *arg.Command) error {
	var client http.Client
	var err error
	var fetchSuccessStatusCode = 200
	var req *http.Request
	var res *http.Response
	var script Script

	var id = "65f10091892c590e57254963" // TODO: Find a way to search through the flags for this value
	var shouldWrite = true              // TODO: make boolean flags work in the parser

	req, err = util.BuildRequest("GET", util.BaseUrl+relativeUrl+"/"+id, nil)
	client = http.Client{}

	res, err = client.Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != fetchSuccessStatusCode {
		return fmt.Errorf("Fetch Status Code %v\n", res.StatusCode)
	}

	err = json.NewDecoder(res.Body).Decode(&script)
	if err != nil {
		return err
	}

	fmt.Printf("ID: %v\nName: %v\nDescription: %v\nLast Modified Date: %v\n", script.Id, script.Name, script.Description, script.LastModified)

	if shouldWrite == true {
		var filename = script.Name + "__" + script.Id + ".js"
		err = os.WriteFile(filename, []byte(script.Content), 0660)
		if err != nil {
			return err
		}
		fmt.Printf("Wrote Contents to file: %v\n", filename)
	}

	return nil
}

func update(_ *arg.Command) error {
	var err error
	var filename = "Test Script 1__65f10091892c590e57254963.js" // TODO: Find a way to get this from the user
	var req *http.Request
	var reqBody string

	var id = parseIdFromFilename(filename)
	if id == "" {
		return fmt.Errorf("Missing Script Id")
	}

	req, err = util.BuildRequest("POST", util.BaseUrl+relativeUrl+"/"+id, reqBody)

	return nil
}
