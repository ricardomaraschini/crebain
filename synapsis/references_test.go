package synapsis

import (
	//	"github.com/davecgh/go-spew/spew"
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
	t.Run("one", func(t *testing.T) {
		// t.Parallel()
		path := cwd + "/internal/fibonacci"
		indexer, err := NewIndexer(path)
		if err != nil {
			t.Fatal(err)
		}

		p, err := indexer.localReferences(path)
		expPackage := Package{
			Name: "github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci",
			Path: path,
			usedSymbols: map[string]struct{}{
				"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci/fib.New":  struct{}{},
				"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci/fib.Till": struct{}{},
			},
		}

		if len(p) != 1 {
			t.Fatal("Unexpected number of packages detected:", len(p))
		}

		checkPackages(t, expPackage, p[0])
	})

	t.Run("two", func(t *testing.T) {
		t.Parallel()
		path := cwd + "/internal/fibonacci"
		pkg2 := path + "/fib"
		indexer, err := NewIndexer(path)
		if err != nil {
			t.Fatal(err)
		}

		p, err := indexer.localReferences(path, pkg2)

		if len(p) != 2 {
			t.Fatal("Unexpected number of packages detected:", len(p))
		}

		want := []Package{
			{
				usedSymbols: map[string]struct{}{},
				Name:        "github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci/fib",
				Path:        pkg2,
			},
			{
				usedSymbols: map[string]struct{}{
					"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci/fib.New":  struct{}{},
					"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci/fib.Till": struct{}{},
				},
				Name: "github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci",
				Path: path,
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
