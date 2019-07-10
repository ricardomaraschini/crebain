package synapsis

import (
	"go/ast"
	"go/parser"
	"go/token"
	"unicode"
)

type Declaration string

// This a block of possible declarations.
const (
	ConstDecl Declaration = "const"
	FuncDecl  Declaration = "func"
	TypeDecl  Declaration = "type"
	VarDecl   Declaration = "var"
)

var (
	parseMode = parser.DeclarationErrors
)

type Symbol struct {
	Declaration Declaration
	Name        string
}

// GetExportedSymbols returns the list of all the symbols contained in the file.
func GetExportedSymbols(path string) ([]Symbol, error) {
	fset := token.NewFileSet()

	f, err := parser.ParseFile(fset, path, nil, parseMode)
	if err != nil {
		return nil, err
	}

	symbols := []Symbol{}
	for _, decl := range f.Decls {
		syms, err := exportedSymbols(decl)
		switch {
		case err != nil:
			return nil, err
		case len(syms) == 0:
			continue
		}

		symbols = append(symbols, syms...)
	}

	return symbols, nil
}

func exportedSymbols(decl ast.Decl) ([]Symbol, error) {
	ss := []Symbol{}
	switch d := decl.(type) {
	case *ast.FuncDecl:
		name := d.Name.String()
		if !isExported(name) {
			return nil, nil
		}
		ss = append(ss, Symbol{
			Name:        d.Name.String(),
			Declaration: FuncDecl,
		})
	case *ast.GenDecl:
		switch d.Tok {
		case token.CONST:
			ss = append(ss, getValueNames(ConstDecl, d.Specs)...)
		case token.VAR:
			ss = append(ss, getValueNames(VarDecl, d.Specs)...)
		case token.TYPE:
			ss = append(ss, getTypesSymbols(d.Specs)...)
		default:
			return nil, nil
		}
	default:
		return nil, nil
	}
	return ss, nil
}

func isExported(name string) bool {
	runes := []rune(name)
	if len(runes) < 1 {
		return false
	}

	return unicode.IsUpper(runes[0])
}

func getValueNames(decl Declaration, specs []ast.Spec) []Symbol {
	ss := []Symbol{}
	for _, spec := range specs {
		names := spec.(*ast.ValueSpec).Names
		for _, n := range names {
			name := n.String()
			if !isExported(name) {
				continue
			}
			ss = append(ss, Symbol{
				Declaration: decl,
				Name:        name,
			})
		}
	}

	return ss
}

func getTypesSymbols(specs []ast.Spec) []Symbol {
	ss := []Symbol{}
	for _, spec := range specs {
		ts := spec.(*ast.TypeSpec)
		name := ts.Name.String()
		if !isExported(name) {
			continue
		}
		ss = append(ss, Symbol{
			Declaration: TypeDecl,
			Name:        name,
		})
	}

	return ss
}
