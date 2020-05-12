package main

/*
#include <stdlib.h>
#include <stdio.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

func Random() int {
	return int(C.random())
}

func Seed(i int) {
	C.srandom(C.uint(i))
}

func Print(s string) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.fputs(cs, (*C.FILE)(C.stdout))
}

func main() {
	Seed(100)
	fmt.Println(Random())
	Print("xxx")
}
