package synapsis

import (
	"reflect"
	"testing"
)

func TestGetExportedSymbols(t *testing.T) {
	path := "./testdata/sample.go"
	expected := []Symbol{
		{
			Declaration: FuncDecl,
			Name:        "ExportedFn",
		},
		{
			Declaration: TypeDecl,
			Name:        "TypeExported",
		},
		{
			Declaration: VarDecl,
			Name:        "VarExported",
		},
		{
			Declaration: ConstDecl,
			Name:        "ConstExported",
		},
	}

	got, err := GetExportedSymbols(path)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(got, expected) {
		t.Fatalf("Wrong symbols returned: %#v", got)
	}
}
