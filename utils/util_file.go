package utils

import (
	"io"
	"os"
)

func ReadFileContent(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read file content using ReadAll
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return content, nil
}
