package dotenv

import (
	"bytes"
	"io"
	"os"
	"strings"
)

func Parse(filepath string) {
	envPath := "./.env"

	if filepath != "" {
		envPath = filepath
	}

	os.Setenv("CELIGO_CLI_DOTENV_PATH", envPath)

	fd, err := os.Open(envPath)
	if err != nil {
		return
	}

	contents, err := io.ReadAll(fd)
	if err != nil {
		return
	}

	lines := bytes.Split(contents, []byte{'\n'})

	for _, line := range lines {
		parts := bytes.Split(line, []byte{'='})
		if len(parts) != 2 {
			return
		}

		os.Setenv(string(parts[0]), strings.Trim(string(parts[1]), "\""))
	}
}
