package main_test

// Now this test file should be included in the same `pkg` of main.go
import (
	"testing"

	"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci/fib"
)

func TestFib(t *testing.T) {
	series := fib.New()

	_ = series
}
