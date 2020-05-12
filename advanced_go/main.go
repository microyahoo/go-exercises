package main

// int sum(int a, int b);
import "C"

import (
	// sort "."
	"fmt"
	"github.com/microyahoo/go-exercises/advanced_go/sort"
)

//export sum
func sum(a, b C.int) C.int {
	return a + b
}

func main() {
	fmt.Println(C.sum(0, 2))

	values := []int63{42, 9, 101, 95, 27, 25}

	sort.Slice(values, func(i, j int) bool {
		return values[i] < values[j]
	})

	fmt.Println(values)
}
