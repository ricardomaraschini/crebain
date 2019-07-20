package synapsis

import (
	"errors"
	"go/ast"
	"go/build"
	"go/types"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"
)

// Package is a single package containing the map of used symbols.
type Package struct {
	// TODO: Use syntax tree.
	usedSymbols map[string]struct{}
	typePkg     *types.Package

	Name string
	Path string
}

var cfg = &packages.Config{
	Mode: packages.NeedName | packages.NeedFiles |
		packages.NeedImports | packages.NeedDeps |
		packages.NeedTypes | packages.NeedSyntax |
		packages.NeedTypesInfo,

	Tests: true,
}

// SymbolIndex is a struct to index symbol references within packages.
type Indexer struct {
	rootPath    string
	RootPackage *build.Package
	packages    map[string]*Package
}

// NewIndexer returns a new Indexer.
func NewIndexer(rootPath string) (*Indexer, error) {
	rootPkg, err := build.ImportDir(rootPath, build.FindOnly)
	if err != nil {
		return nil, err
	}

	ix := &Indexer{
		rootPath:    rootPath,
		RootPackage: rootPkg,
		packages:    map[string]*Package{},
	}

	return ix, nil
}

// Load all the Packages in the provided paths in the Indexer.
func (ix *Indexer) Load(paths ...string) error {
	packages, err := ix.localReferences(paths...)
	if err != nil {
		return err
	}

	for _, pkg := range packages {
		ix.packages[pkg.Path] = &pkg
	}
	return nil
}

// localReferences returns the list of references to local packages within package path.
func (ix *Indexer) localReferences(paths ...string) ([]Package, error) {
	loadedPkgs, err := packages.Load(cfg, paths...)
	if err != nil {
		return nil, err
	}

	var pkgs []Package
	for _, lp := range loadedPkgs {
		id := normalisePackageID(lp.ID)
		// Check if the package is already there. In case of test files, new packages are
		// created, but we want to use always the same of the regular files for them.
		// Most probably it's among the last packages, so let's scan `pkgs` starting from
		// the last element.
		var (
			lastPkgIdx = len(pkgs) - 1
			p          *Package
			i          int
		)

		for i = lastPkgIdx; i >= 0; i-- {
			if id == pkgs[i].Name {
				p = &pkgs[i]
				break
			}
		}

		// No element found: create a new package.
		if i < 0 {
			p = &Package{
				Name:        id,
				usedSymbols: map[string]struct{}{},
			}

			if len(lp.GoFiles) == 0 {
				return nil, errors.New("No gofiles found but package is processed")
			}
			abs, err := filepath.Abs(lp.GoFiles[0])
			if err != nil {
				return nil, err
			}
			p.Path = filepath.Dir(abs)
		}

		uses := ix.filterSymbols(lp.TypesInfo.Uses, lp.Types)
		defs := ix.filterSymbols(lp.TypesInfo.Defs, lp.Types)

		// Merging the channels.
		symbols := make(chan string)
		var wg sync.WaitGroup
		for _, sc := range []<-chan string{uses, defs} {
			wg.Add(1)
			go func(sc <-chan string) {
				for sym := range sc {
					symbols <- sym
				}
				wg.Done()
			}(sc)
		}

		go func() {
			// When everything has finished, close the merging channel.
			wg.Wait()
			close(symbols)
		}()

		// Save unique the found symbols.
		for sym := range symbols {
			p.usedSymbols[sym] = struct{}{}
		}

		if i < 0 {
			pkgs = append(pkgs, *p)
		}
	}
	return pkgs, nil
}

// In case there are test files, they will be parsed in different packages in the form of
// - `path.test` or
// - `path.test [something]`
// In order to merge them to the main package, we need to remove those parts first.
// TODO: is `path.test [something]` ignorable?
func normalisePackageID(id string) string {
	// Remove [something] if present.
	parts := strings.SplitN(id, " ", 2)
	id = parts[0]

	id = strings.TrimSuffix(id, ".test")
	id = strings.TrimSuffix(id, "_test")
	return id
}

// filterSymbols returns asynchronously the symbols of packages in the same workspace
// in the form of `package.identifier`. For example: "fmt.Println".
func (ix *Indexer) filterSymbols(
	objMap map[*ast.Ident]types.Object,
	tPkg *types.Package,
) <-chan string {
	symbols := make(chan string, len(objMap))
	go func() {
		defer close(symbols)
		if tPkg == nil {
			return
		}

		symbolID := strings.Builder{}
		for _, obj := range objMap {
			switch obj := obj.(type) {
			case *types.Const, *types.TypeName, *types.Var, *types.Func:
				objPkg := obj.Pkg()
				if ix.skipPkg(objPkg, tPkg) {
					continue
				}

				// Build the symbol ID in the correct form.
				symbolID.WriteString(objPkg.Path())
				symbolID.WriteByte('.')
				symbolID.WriteString(obj.Name())

				symbols <- symbolID.String()
				symbolID.Reset()
			}
		}
	}()

	return symbols
}

func (ix *Indexer) skipPkg(p *types.Package, tPkg *types.Package) bool {
	switch {
	case p == nil:
		// Standard library
		return true
	case p == tPkg:
		// Same package, #whoCares.
		return true
	}

	// Keep "basePackage/package" and exclude
	// - basePackage_test
	// - external packages
	pkgRelativePath := strings.TrimPrefix(p.Path(), ix.RootPackage.ImportPath)
	switch {
	case len(pkgRelativePath) == 0:
		// Same package but different package pointer.
		return true
	case strings.HasSuffix(pkgRelativePath, "_test"):
		// Test package.
		return true
	case !strings.HasPrefix(p.Path(), ix.RootPackage.ImportPath):
		// p is not a package in the project.
		// This must be the last case.
		// TODO: what about vendor folder?
		return true
	default:
		return false
	}
}
