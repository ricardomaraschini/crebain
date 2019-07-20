package trunner

import (
	"bytes"
	"fmt"
	"testing"
)

func TestInvalidDir(t *testing.T) {
	r := New()
	res, err := r.Run("/tmp/does-not-exist")
	if res == nil {
		t.Fatal("expected result for invalid test directory")
	}
	if err != nil {
		t.Fatal("test on invalid dir returing error", err)
	}

	if res.Code == 0 {
		t.Fatal("expected result code not be be success(0)")
	}

	if len(res.Out) != 1 {
		t.Fatal("expected only one output on the result")
	}
}

type readerCrash struct{}

func (r *readerCrash) Read([]byte) (int, error) {
	return 0, fmt.Errorf("reader crash error")
}

func TestParseOutput(t *testing.T) {
	r := New()

	content := bytes.NewBuffer([]byte("1\n2\n\n3\n"))
	lines, err := r.parseOutput(content)
	if err != nil {
		t.Fatal("unexpected error parsing output")
	}
	if len(lines) != 3 {
		t.Fatal("invalid number of lines returned by parseOutput")
	}

	_, err = r.parseOutput(&readerCrash{})
	if err == nil {
		t.Fatal("error expected but nil received")
	}
	if err.Error() != "reader crash error" {
		t.Fatalf("unexpected error %q", err)
	}
}
