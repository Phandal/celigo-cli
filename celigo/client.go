package celigo

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const baseUrl = "https://api.integrator.io/v1"
const envFile = ".env"
const celigoCliEnvKey = "CELIGO_API_KEY"

var memoizedApiKey string

type request struct {
	Method    string
	Url       string
	Body      *bytes.Buffer
	Code      int
	Resources any
}

type response struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func apiKey() (string, error) {
	apiKey := os.Getenv("CELIGO_API_KEY")

	if len(apiKey) == 0 {
		return "", errors.New("missing CELIGO_API_KEY")
	}

	return apiKey, nil
}

func buildRequest(method string, url string, body *bytes.Buffer) (*http.Request, error) {
	var err error
	var req *http.Request
	apikey, err := apiKey()
	if err != nil {
		return nil, err
	}

	if body != nil {
		req, err = http.NewRequest(method, url, body)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+apikey)
	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	req.Header.Add("Accept", "application/json")

	return req, nil
}

func callApi(req *http.Request) (*http.Response, error) {
	var client http.Client
	client = http.Client{}

	return client.Do(req)
}

func checkStatusCode(res *http.Response, code int) bool {
	return res.StatusCode == code
}

func decodeResponse(input io.ReadCloser, parsedRecord any) error {
	return json.NewDecoder(input).Decode(&parsedRecord)
}

func encodeBody(resource any) ([]byte, error) {
	return json.Marshal(resource)
}

func executeRequest(celigoRequest *request) error {
	var req *http.Request
	var res *http.Response
	var err error

	if req, err = buildRequest(celigoRequest.Method, celigoRequest.Url, celigoRequest.Body); err != nil {
		return err
	}

	if res, err = callApi(req); err != nil {
		return err
	}

	if !checkStatusCode(res, celigoRequest.Code) {
		return fmt.Errorf("Unsuccessful Request. Status Code: %d", res.StatusCode)
	}

	if decodeResponse(res.Body, celigoRequest.Resources) != nil {
		return err
	}

	return nil
}

func newCeligoRequest(method string, relativeUrl string, body *bytes.Buffer, code int, resources any) *request {
	return &request{
		Method:    method,
		Url:       baseUrl + relativeUrl,
		Body:      body,
		Code:      code,
		Resources: resources,
	}
}

func ExecuteGet(relativeUrl string, code int, returnResource any) error {
	var req = newCeligoRequest("GET", relativeUrl, nil, code, returnResource)
	return executeRequest(req)
}

func ExecutePost(relativeUrl string, resource any, code int, returnResource any) error {
	var content = []byte{}
	var err error

	if content, err = encodeBody(resource); err != nil {
		return err
	}

	var req = newCeligoRequest("POST", relativeUrl, bytes.NewBuffer(content), code, returnResource)
	return executeRequest(req)
}

func ExecutePut(relativeUrl string, resource any, code int, returnResource any) error {
	var content = []byte{}
	var err error

	if content, err = encodeBody(resource); err != nil {
		return err
	}

	var req = newCeligoRequest("PUT", relativeUrl, bytes.NewBuffer(content), code, returnResource)
	return executeRequest(req)
}

func ExecuteDelete(relativeUrl string, code int) error {
	req := newCeligoRequest("DELETE", relativeUrl, nil, code, nil)
	return executeRequest(req)
}
