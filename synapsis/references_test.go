package synapsis

import (
	"os"
	"reflect"
	"testing"
)

func TestLocalReferences(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	path := cwd + "/internal/fibonacci"

	p, err := localReferences(path)
	expectedSymbols := map[string]struct{}{
		"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci/fib.New":  struct{}{},
		"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci/fib.Till": struct{}{},
	}
	switch {
	case len(p) != 1:
		t.Fatal("Unexpected number of packages detected:", len(p))
	case !reflect.DeepEqual(p[0].usedSymbols, expectedSymbols):
		t.Fatalf("Symbols don't match expected: %#v", p[0].usedSymbols)
	}
}