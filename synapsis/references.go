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

	// TODO: Make a test file in internal/fibonacci.
	Tests: true,
}

// localReferences returns the hashtable of references to local packages within packages in paths.
func localReferences(paths ...string) ([]pkg, error) {
	loadedPkgs, err := packages.Load(cfg, paths...)
	if err != nil {
		return nil, err
	}

	var pkgs []pkg
	for _, lp := range loadedPkgs {
		p := pkg{
			name:        lp.ID,
			pPkg:        lp.Types,
			usedSymbols: map[string]struct{}{},
		}

		uses := p.filterSymbols(lp.TypesInfo.Uses)
		defs := p.filterSymbols(lp.TypesInfo.Defs)

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

		pkgs = append(pkgs, p)
	}
	return pkgs, nil
}

// filterSymbols returns asynchronously the symbols in the form of `package.identifier`.
// For example: "fmt.Println".
func (p *pkg) filterSymbols(objMap map[*ast.Ident]types.Object) <-chan string {
	symbols := make(chan string, len(objMap))
	go func() {
		defer close(symbols)
		if p.pPkg == nil {
			return
		}

		symbolID := strings.Builder{}
		for _, obj := range objMap {
			switch obj := obj.(type) {
			case *types.Const, *types.TypeName, *types.Var, *types.Func:
				objPkg := obj.Pkg()
				if p.skipPkg(objPkg) {
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

func (bp *pkg) skipPkg(p *types.Package) bool {
	switch {
	case p == nil:
		// Standard library
		return true
	case p == bp.pPkg:
		// Same package, #whoCares.
		return true
	case !strings.HasPrefix(p.Path(), bp.pPkg.Path()):
		// p is not a package in the project.
		// TODO: use crebain's `path` package.
		return true
	default:
		return false
	}
}
