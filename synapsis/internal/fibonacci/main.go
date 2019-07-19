package main

import (
	_ "errors"
	"fmt"
	. "os"
	"strconv"

	"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci/fib"
)

func main() {
	if len(Args) < 2 {
		fmt.Println("Number not set")
	}

	number, err := strconv.Atoi(Args[1])
	if err != nil {
		fmt.Println("not valid number")
		Exit(1)
	}

	series := fib.New()

	fmt.Println("Asking", number)
	result := series.Till(uint(number))
	fmt.Println("Result is", result)
}
