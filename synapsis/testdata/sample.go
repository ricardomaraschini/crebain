// Package one is a package.
package one

import "fmt"

func unexportedFn() {
	fmt.Println("I can't remember anything")
}

// ExportedFn ...
func ExportedFn() {
	fmt.Println("Can't tell if this is true or dream")
}

// Landmine has taken my sight.

type unexportedType string

// TypeExported ...
type TypeExported struct {
	ohPleaseGodHelpMe int
}

var varUnexported = "Deep down inside I feel to scream"

// VarExported ...
var VarExported = "This terrible silence stops me"

// Taken my speech.
// Taken my hearing.
// Taken my arms.

// Taken my legs
const (
	constUnexported = "Hold my breath as I wish for death"
	// ConstExported has Taken my soul.
	ConstExported = "Oh please, God, wake me"
)

const (
	a = 12
	b = 234
)

// Left me with life in hell
