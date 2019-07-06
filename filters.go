package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// hasTestFile is a filter that returns true if either filePath points to a
// go test file or points to a go file with correspondent go test file.
func hasTestFile(filePath string) bool {
	if filepath.Ext(filePath) != ".go" {
		return false
	}

	if strings.HasSuffix(filePath, "_test.go") {
		return true
	}

	ext := filepath.Ext(filePath)
	testFilePath := fmt.Sprintf(
		"%s_test.go",
		filePath[0:len(filePath)-len(ext)],
	)
	_, err := os.Stat(testFilePath)
	return err == nil
}
