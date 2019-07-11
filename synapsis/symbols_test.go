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
			Start:       147,
			End:         157,
		},
		{
			Declaration: TypeDecl,
			Name:        "TypeExported",
			Start:       303,
			End:         315,
		},
		{
			Declaration: VarDecl,
			Name:        "VarExported",
			Start:       431,
			End:         442,
		},
		{
			Declaration: ConstDecl,
			Name:        "ConstExported",
			Start:       658,
			End:         671,
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
