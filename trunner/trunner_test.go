package trunner

import (
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
