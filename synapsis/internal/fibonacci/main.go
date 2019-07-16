package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/ricardomaraschini/crebain/synapsis/internal/fibonacci/fib"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Number not set")
	}

	number, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("not valid number")
		os.Exit(1)
	}

	series := fib.New()

	fmt.Println("Asking", number)
	result := series.Till(uint(number))
	fmt.Println("Result is", result)
}
