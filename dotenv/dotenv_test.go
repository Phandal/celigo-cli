package dotenv

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	Parse("./.env_test")
	envFile := os.Getenv("CELIGO_CLI_DOTENV_PATH")

	if envFile != "./.env_test" {
		t.Fatalf("Failed to set CELIGO_CLI_DOTENV_PATH")
	}

	if os.Getenv("TEST_KEY") != "new_value" {
		t.Fatalf("Failed to set environment variables with dotenv")
	}

	if os.Getenv("QUOTE_KEY") != "quote value" {
		t.Fatalf("Failed to set quoted environment variables with dotenv")
	}
}
