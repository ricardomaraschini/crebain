package trunner

import (
	"bufio"
	"bytes"
	"encoding/json"
	"os/exec"
	"time"
)

// New returns a new TRunner.
func New() *TRunner {
	return &TRunner{}
}

// TestResult holds a result of a test execution.
type TestResult struct {
	Code int
	Out  []ResultLine
}

// ResultLine holds every line of a go test.
type ResultLine struct {
	Time    time.Time
	Action  string
	Package string
	Test    string
	Output  string
	Elapsed float64
}

// TRunner is go test helper.
type TRunner struct{}

// Run runs tests on provided directories.
func (t *TRunner) Run(dir string) (*TestResult, error) {
	var result TestResult

	stdout := bytes.NewBuffer(nil)
	stderr := bytes.NewBuffer(nil)
	cmd := exec.Command("go", "test", "-cover", "-json", dir)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
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

	out, err := t.parseTestOutput(stdout)
	if err != nil {
		return nil, err
	}
	result.Out = out
	return &result, nil
}

// Parses the buffer, unmarshals all line into ResultLine structs.
func (t *TRunner) parseTestOutput(buf *bytes.Buffer) ([]ResultLine, error) {
	var lines []ResultLine

	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		lres := ResultLine{}
		if err := json.Unmarshal(scanner.Bytes(), &lres); err != nil {
			return nil, err
		}
		lines = append(lines, lres)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}
