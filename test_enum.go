package main

import "fmt"

type taskOp uint32

const (
	A taskOp = iota
	B
	C
)

func (o taskOp) String() string {
	var operation string
	switch o {
	case A:
		operation = "A"
	case B:
		operation = "B"
	case C:
		operation = "C"
	default:
		operation = "unknown"
	}
	return operation
}

func main() {
	var x, y taskOp
	x = 2
	fmt.Println(x)
	y = 7
	fmt.Println(y)
}
