package util

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

const BaseUrl = "https://api.integrator.io/v1"

var apiKey string

func ApiKey() string {
	if apiKey == "" {
		apiKey = os.Getenv("CELIGO_API_KEY")
	}

	return apiKey
}

func BuildRequest(method string, url string, body io.Reader) (*http.Request, error) {
	var err error
	var req *http.Request

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+ApiKey())

	return req, nil
}

func CheckStatusCode(res *http.Response, code int) bool {
	return res.StatusCode == code
}

func DecodeResponse(input io.ReadCloser, parsedRecord interface{}) error {
	return json.NewDecoder(input).Decode(&parsedRecord)
}

func ExecuteRequest(req *http.Request) (*http.Response, error) {
	var client http.Client
	client = http.Client{}

	return client.Do(req)
}
