// https://medium.com/a-journey-with-go/go-how-does-the-garbage-collector-mark-the-memory-72cfc12c6976
package main

import (
	"fmt"
	"runtime"
)

type struct1 struct {
	a, b int64
	c, d float64
	e    *struct2
}

type struct2 struct {
	f, g int64
	h, i float64
}

func main() {
	s1 := allocStruct1()
	s2 := allocStruct2()

	func() {
		_ = allocStruct2()
	}()

	runtime.GC()

	fmt.Printf("s1 = %X, s2 = %X\n", &s1, &s2)
}

//go:noinline
func allocStruct1() *struct1 {
	return &struct1{
		e: allocStruct2(),
	}
}

//go:noinline
func allocStruct2() *struct2 {
	return &struct2{}
}
