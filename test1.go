package main

import (
	"fmt"
)

func main() {
	fmt.Println("hello world!")

	const PtrSize = 4 << (^uintptr(0) >> 63)
	fmt.Printf("(^uintptr(0))=%b\n", (^uintptr(0)))
	fmt.Printf("(^uintptr(0) >> 63)=%d\n", (^uintptr(0) >> 63))
	fmt.Printf("ptrSize=0b%b\n", PtrSize)
}
