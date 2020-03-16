package main

import (
	"fmt"
)

func main() {
	fmt.Println("First round: without fallthrough")
	switch 1 {
	case 0:
		fmt.Println(0)
	case 1:
		fmt.Println(1)
	case 2:
		fmt.Println(2)
	case 3:
		fmt.Println(3)
	}

	fmt.Println("Second round: with fallthrough")
	switch 1 {
	case 0:
		fmt.Println(0)
		fallthrough
	case 1:
		fmt.Println(1)
		fallthrough
	case 2:
		fmt.Println(2)
		fallthrough
	case 3:
		fmt.Println(3)
		fallthrough
	default:
		fmt.Println("default")
		// cannot fallthrough final case in switch
		// fallthrough
	}
}
