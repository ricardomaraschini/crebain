package synapsis

import (
	"go/ast"
	"go/types"
	"strings"
	"sync"

	"golang.org/x/tools/go/packages"
)

type pkg struct {
	name string
	// TODO: Use syntax tree.
	usedSymbols map[string]struct{}
	pPkg        *types.Package
}

var cfg = &packages.Config{
	// TODO: Take only what you need.
	Mode: packages.NeedName | packages.NeedFiles |
		packages.NeedCompiledGoFiles | packages.NeedImports |
		packages.NeedDeps | packages.NeedExportsFile |
		packages.NeedTypes | packages.NeedSyntax |
		packages.NeedTypesInfo | packages.NeedTypesSizes,

	Tests: true,
}

// localReferences returns the list of references to local packages within package path.
// TODO: use crebain base path.
func localReferences(paths ...string) ([]pkg, error) {
	loadedPkgs, err := packages.Load(cfg, paths...)
	if err != nil {
		return nil, err
	}

	var pkgs []pkg
	for _, lp := range loadedPkgs {
		id := normalisePackageID(lp.ID)
		// Check if package is already there. In case of test files, new packages are
		// created, but we want to use always the same `pkg` for them.
		// Most probably is among the last packages, so let's scan the current `pkgs`
		// starting from the last element.

		var (
			lastPkgIdx = len(pkgs) - 1
			p          *pkg
			i          int
		)

		for i = lastPkgIdx; i >= 0; i-- {
			if id == pkgs[i].name {
				p = &pkgs[i]
				break
			}
		}

		// No element found: create a new package.
		if i < 0 {
			p = &pkg{
				name:        id,
				usedSymbols: map[string]struct{}{},
			}
		}

		uses := p.filterSymbols(lp.TypesInfo.Uses, lp.Types)
		defs := p.filterSymbols(lp.TypesInfo.Defs, lp.Types)

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

// In case there are test file, they will be parsed in different packages in the form of
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
func (p *pkg) filterSymbols(
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
				if p.skipPkg(objPkg, tPkg) {
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

func (bp *pkg) skipPkg(p *types.Package, tPkg *types.Package) bool {
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
	// TODO: use crebain's `path` package.
	pkgRelativePath := strings.TrimPrefix(p.Path(), bp.name)
	switch {
	case len(pkgRelativePath) == 0:
		// Same package but different package pointer.
		return true
	case strings.HasSuffix(pkgRelativePath, "_test"):
		// Test package.
		return true
	case !strings.HasPrefix(p.Path(), bp.name):
		// p is not a package in the project.
		// This must be the last case.
		// TODO: what about vendor folder?
		return true
	default:
		return false
	}
}
