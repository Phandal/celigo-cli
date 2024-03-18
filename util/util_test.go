package util_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Phandal/celigo-cli/util"
)

func TestCheckStatusCode(t *testing.T) {
	var expected = true
	var actual bool
	var successCode = 200
	var res *http.Response
	var err error

	res, err = http.Get("https://dummyjson.com")
	if err != nil {
		t.Errorf("Failed to execute Get request: %v\n", err)
	}

	actual = util.CheckStatusCode(res, successCode)

	if expected != actual {
		t.Fatalf("Expected: %v\tActual: %v\n", expected, actual)
	}
}

func TestCheckStatusCodeInvalid(t *testing.T) {
	var expected = false
	var actual bool
	var successCode = 500
	var res *http.Response
	var err error

	res, err = http.Get("https://dummyjson.com")
	if err != nil {
		t.Errorf("Failed to execute Get request: %v\n", err)
	}

	actual = util.CheckStatusCode(res, successCode)

	if expected != actual {
		t.Fatalf("Expected: %v\tActual: %v\n", expected, actual)
	}
}

func TestDecodeResponse(t *testing.T) {
	type TestResponse struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}

	var input io.ReadCloser
	var actual TestResponse
	var expected TestResponse
	var jsonString string
	var err error

	expected = TestResponse{
		Id:   12345,
		Name: "Test Script",
	}

	jsonString = "{\"id\": 12345,\"name\": \"Test Script\"}"
	input = io.NopCloser(strings.NewReader(jsonString))

	if err = util.DecodeResponse(input, &actual); err != nil {
		t.Errorf("Failed to Decode JSON string: %v\n", err)
	}

	if expected != actual {
		t.Fatalf("Expected: %v\nActual: %v\n", expected, actual)
	}
}
