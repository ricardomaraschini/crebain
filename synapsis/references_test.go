package synapsis

import (
	"github.com/davecgh/go-spew/spew"
	"os"
	"reflect"
	"testing"
)

var (
	cwd string
)

func TestMain(m *testing.M) {
	var err error
	cwd, err = os.Getwd()
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func TestLocalReferences(t *testing.T) {
	t.Run("one pkg", func(t *testing.T) {
		t.Parallel()
		path := cwd + "/internal/fibonacci"
		indexer, err := NewIndexer(path)
		if err != nil {
			t.Fatal(err)
		}

		p, err := indexer.localReferences(path)
		const base = "github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci"
		expPackage := Package{
			Name: base,
			Path: path,
			usedSymbols: map[string]struct{}{
				base + "/fib.New":  struct{}{},
				base + "/fib.Till": struct{}{},
			},
		}

		if len(p) != 1 {
			t.Fatal("Unexpected number of packages detected:", len(p))
		}

		checkPackages(t, expPackage, p[0])
	})

	t.Run("more pkgs", func(t *testing.T) {
		t.Parallel()
		var (
			path = cwd + "/internal/fibonacci"
			pkgs = []string{
				path,
				path + "/fib",
				path + "/fib/esoteric",
				path + "/useless",
			}
		)

		indexer, err := NewIndexer(path)
		if err != nil {
			t.Fatal(err)
		}

		p, err := indexer.localReferences(pkgs...)

		if len(p) != 4 {
			t.Fatal("Unexpected number of packages detected:", len(p))
		}

		spew.Dump(p)
		t.Fatal("This will be my monument")
		const base = "github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci"
		want := []Package{
			{
				usedSymbols: map[string]struct{}{
					// there should be something here.
				},
				Name: base + "/fib/esoteric",
				Path: pkgs[2],
			},
			{
				usedSymbols: map[string]struct{}{
					base + "/fib.Sequence":             struct{}{}, // this doesn't make any sense.
					base + "/fib/esoteric.SuperFastDT": struct{}{},
				},
				Name: base + "/fib/useless",
				Path: pkgs[3],
			},
			{
				usedSymbols: map[string]struct{}{},
				Name:        base + "/fib",
				Path:        pkgs[1],
			},
			{
				usedSymbols: map[string]struct{}{
					base + "/fib.New":  struct{}{},
					base + "/fib.Till": struct{}{},
				},
				Name: base,
				Path: pkgs[0],
			},
		}
		for i, w := range want {
			checkPackages(t, w, p[i])
		}
	})
}

func checkPackages(t *testing.T, want, got Package) {
	switch {
	case want.Name != got.Name:
		t.Errorf("Wrong name. Want: %s, Got: %s", want.Name, got.Name)
	case want.Path != got.Path:
		t.Errorf("Wrong path. Want: %s, Got: %s", want.Path, got.Path)
	case !reflect.DeepEqual(want.usedSymbols, got.usedSymbols):
		t.Errorf("Symbols don't match want: %+v", got.usedSymbols)
	}

}

func TestLoad(t *testing.T) {
	path := cwd + "/internal/fibonacci"
	indexer, err := NewIndexer(path)
	if err != nil {
		t.Fatal(err)
	}

	if err := indexer.Load(path, path+"/fib"); err != nil {
		t.Fatal(err)
	}

	expKeys := []string{
		path, path + "/fib",
	}

	for _, k := range expKeys {
		if _, ok := indexer.packages[k]; !ok {
			//spew.Dump(indexer.packages)
			t.Fatalf("Key %s not loaded", k)
		}
	}
}

func TestNormaliseImportPath(t *testing.T) {
	testCases := []string{
		"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci",
		"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci.test",
		"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci [github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci.test]",
		"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci_test [github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci.test]",
	}
	exp := "github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci"

	for _, id := range testCases {
		got := normaliseImportPath(id)
		if got != exp {
			t.Fatalf("Not normalised correctly: %q", got)
		}
	}
}

func BenchmarkNormaliseImportPath(b *testing.B) {
	testCases := []string{
		"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci",
		"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci.test",
		"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci [github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci.test]",
		"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci_test [github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci.test]",
	}

	for _, tc := range testCases {
		b.Run(tc, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				normaliseImportPath(tc)
			}
		})
	}
}

func BenchmarkLoadPackages(b *testing.B) {
	var (
		pkg1 = cwd + "/internal/fibonacci"
		pkg2 = cwd + "/internal/fibonacci/fib"
	)

	indexer, err := NewIndexer(cwd)
	if err != nil {
		b.Fatal(err)
	}

	for n := 0; n < b.N; n++ {
		_, _ = indexer.localReferences(pkg1, pkg2)
	}
}
