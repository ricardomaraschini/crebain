package fib

import "testing"

func TestTill(t *testing.T) {
	seq := New()

	got := seq.Till(6)
	exp := 13
	if got != exp {
		t.Fatal("It didn't work", got)
	}
}
