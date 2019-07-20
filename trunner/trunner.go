package trunner

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
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

// Content return the result of the test.
func (t *TestResult) Content() []string {
	return t.Out
}

// Title returns the test directory to be rendered as test title.
func (t *TestResult) Title() string {
	return t.Dir
}

// Success return if the test was successful.
func (t *TestResult) Success() bool {
	return t.Code == 0
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

	lines, err := t.parseOutput(std)
	if err != nil {
		return nil, err
	}

	result.Out = lines
	return &result, nil
}

// parseOutput splits content by new lines, removing empty lines.
func (t *TRunner) parseOutput(content io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(content)
	lines := make([]string, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		lines = append(lines, line)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
