package main

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

func main() {
	_ = sync.Cond{}
	fmt.Println("vim-go")

	var a strings.Builder
	a.Write([]byte("a"))
	// b := a
	// b.Write([]byte("b")) // panic: strings: illegal use of non-zero Builder copied by value

	var m cond
	m.checker.check()
	n := m
	n.checker.check()
}

type cond struct {
	checker copyChecker
}

type copyChecker uintptr

func (c *copyChecker) check() {
	fmt.Printf("Before: c: %12v, *c: %12v, uintptr(*c): %12v, uintptr(unsafe.Pointer(c)): %12v\n",
		c, *c, uintptr(*c), uintptr(unsafe.Pointer(c)))
	swapped := atomic.CompareAndSwapUintptr((*uintptr)(c), 0, uintptr(unsafe.Pointer(c)))
	fmt.Printf("After : c: %12v, *c: %12v, uintptr(*c): %12v, uintptr(unsafe.Pointer(c)): %12v, swapped: %5v\n",
		c, *c, uintptr(*c), uintptr(unsafe.Pointer(c)), swapped)
}

// Before: c: 0xc0000180b8, *c:            0, uintptr(*c):            0, uintptr(unsafe.Pointer(c)): 824633819320
// After : c: 0xc0000180b8, *c: 824633819320, uintptr(*c): 824633819320, uintptr(unsafe.Pointer(c)): 824633819320, swapped:  true
// Before: c: 0xc0000180f0, *c: 824633819320, uintptr(*c): 824633819320, uintptr(unsafe.Pointer(c)): 824633819376
// After : c: 0xc0000180f0, *c: 824633819320, uintptr(*c): 824633819320, uintptr(unsafe.Pointer(c)): 824633819376, swapped: false
