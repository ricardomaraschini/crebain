package useless

import (
	"fmt"
	"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci/fib"
	"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci/fib/esoteric"
)

func IDoNothing() {
	fmt.Println("NOTHING NOTHING NOTHING")
	fmt.Println(esoteric.SuperFastDT("N"))
}

type CanISeeIt func(fib.Sequence) error

type Bounty int
