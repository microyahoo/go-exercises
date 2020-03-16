package main

import (
	"fmt"
)

// A is a test struct
type A struct {
	aI int
	aS string
	aB *B
}

// B is a test
type B struct {
	bS   string
	bMap map[string]int
}

func main() {
	b := &B{
		bS: "abc",
		bMap: map[string]int{
			"a": 1,
			"b": 2,
		},
	}
	a := &A{
		aI: 12,
		aS: "hello world",
		aB: b,
	}
	fmt.Printf("The address of a: %p, a=%v\n", a, a)
	fmt.Printf("The address of b: %p, b=%v\n", b, b)

	x := *a
	fmt.Printf("The address of x: %p, x=%v\n", &x, x)
	x.aB.bMap["c"] = 3
	fmt.Printf("\nThe address of a: %p, a=%v\n", a, a)
	fmt.Printf("The address of b: %p, b=%v\n", b, b)
	fmt.Printf("The address of x: %p, x=%v\n", &x, x)

}

// The address of a: 0xc00000c0a0, a=&{12 hello world 0xc00000c080}
// The address of b: 0xc00000c080, b=&{abc map[a:1 b:2]}
// The address of x: 0xc00000c0e0, x={12 hello world 0xc00000c080}
