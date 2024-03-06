package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const API_KEY = "CELIGO_API_KEY"

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

func usage() {
	fmt.Printf("usage: celigo {action}\n")
}

func getApiKey() (apikey string, err any) {
	key := os.Getenv(API_KEY)
	if len(key) == 0 {
		return "", "Empty API Key"
	} else {
		return key, nil
	}
}

func fetch(apiKey string) {
	request, err := http.NewRequest("GET", "https://api.integrator.io/v1/scripts/65e67a82c1774f84e5c27304", nil)
	if err != nil {
		fmt.Printf("Error creating new request: %s\n", err)
		return
	}

	request.Header.Add("Authorization", "Bearer "+apiKey)

	client := http.Client{}

	res, err := client.Do(request)
	if err != nil {
		fmt.Printf("Error making request: %s\n", err)
		return
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}

	fmt.Printf("Raw Body:\n%s\n", body)
	var script Script
	if err = json.Unmarshal(body, &script); err != nil {
		fmt.Printf("Error parsing json data: %s\n", err)
		return
	}

	fmt.Printf("Script:\nName: %s\nDescription: %s\nId: %s\nContents:\n", script.Name, script.Description, script.Id)
	fmt.Println(script.Content)

	file, err := os.Create(script.Name + ".js")
	if err != nil {
		fmt.Printf("Error creating file: %s\n", err)
		return
	}

	file.WriteString(script.Content)
	file.Close()
	fmt.Printf("Wrote to file %s\n", file.Name())
}

func update(apiKey string) {
	contents, err := os.ReadFile("Test Script.js")
	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	var updateScript Script
	updateScript.Id = "65e67a82c1774f84e5c27304"
	updateScript.Name = "Test Script Update"
	updateScript.Description = "Test Script Description"
	updateScript.Content = bytes.NewBuffer(contents).String()

	reqBody, err := json.Marshal(updateScript)
	if err != nil {
		fmt.Printf("Error marshalling json: %s\n", err)
		return
	}

	request, err := http.NewRequest(
		"PUT",
		"https://api.integrator.io/v1/scripts/65e67a82c1774f84e5c27304",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		fmt.Printf("Error creating new request: %s\n", err)
		return
	}

	request.Header.Add("Authorization", "Bearer "+apiKey)
	request.Header.Add("Content-Type", "application/json")
	client := http.Client{}

	res, err := client.Do(request)
	if err != nil {
		fmt.Printf("Error making request: %s\n", err)
		return
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		fmt.Printf("Error reading response body: %s\n", err)
		return
	}

	fmt.Printf("Raw Body:\n%s\n", body)
	var script Script
	if err = json.Unmarshal(body, &script); err != nil {
		fmt.Printf("Error parsing json data: %s\n", err)
		return
	}

	fmt.Printf("Script:\nName: %s\nDescription: %s\nId: %s\nContents:\n", script.Name, script.Description, script.Id)
	fmt.Println(script.Content)
}

func main() {
	if len(os.Args) != 2 {
		usage()
		return
	}

	action := strings.ToLower(os.Args[1])
	apiKey, err := getApiKey()
	if err != nil {
		fmt.Printf("Error reading env api key: %s\n", err)
		return
	}

	switch action {
	case "fetch":
		fetch(apiKey)
	case "update":
		update(apiKey)
	}
}
