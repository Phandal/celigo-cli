package util_test

import (
	"net/http"
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
