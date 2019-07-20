package trunner

import (
	"bytes"
	"io/ioutil"
	"os/exec"
	"strings"
)

// New returns a new TRunner.
func New() *TRunner {
	return &TRunner{}
}

// TestResult holds a result of a test execution.
type TestResult struct {
	Out  []string
	Dir  string
	Code int
}

// TRunner is go test helper.
type TRunner struct{}

// Run runs tests on provided directories.
func (t *TRunner) Run(dir string) (*TestResult, error) {
	result := TestResult{
		Dir: dir,
	}

	std := bytes.NewBuffer(nil)
	cmd := exec.Command("go", "test", "-cover", dir)
	cmd.Stdout = std
	cmd.Stderr = std
	if err := cmd.Start(); err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		// on a failure scenario, try to capture the exit code.
		exiterr, ok := err.(*exec.ExitError)
		if !ok {
			return nil, err
		}
		result.Code = exiterr.ExitCode()
	}

	rawTestOutput, err := ioutil.ReadAll(std)
	if err != nil {
		return nil, err
	}

	result.Out = t.parseOutput(string(rawTestOutput))
	return &result, nil
}

// parseOutput splits content by new lines, removing empty lines.
func (t *TRunner) parseOutput(content string) []string {
	rawLines := strings.Split(content, "\n")
	lines := make([]string, 0)
	for _, line := range rawLines {
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	return lines
}
