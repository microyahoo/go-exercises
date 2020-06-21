package main

import (
	"fmt"
	"runtime"
	"unsafe"
)

type funcTwiceClosure struct {
	f uintptr
	x int
}

func ptrToFunc(p unsafe.Pointer) func() int

func asmFuncTwiceClosureAddr() uintptr

func asmFuncTwiceClosureBody() int

func syscallWrite_darwin(fd int, msg string) int

func newTwiceFuncClosure(x int) func() int {
	var p = &funcTwiceClosure{
		f: asmFuncTwiceClosureAddr(),
		x: x,
	}
	return ptrToFunc(unsafe.Pointer(p))
}

func main() {
	fnTwice := newTwiceFuncClosure(2)
	fmt.Println(fnTwice())
	fmt.Println(fnTwice())
	fmt.Println(fnTwice())

	if runtime.GOOS == "darwin" {
		syscallWrite_darwin(1, "hello world....")
	}
}
