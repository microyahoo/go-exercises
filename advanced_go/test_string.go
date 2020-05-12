package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	s := "hello world"

	sHdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	fmt.Println(sHdr)
	fmt.Println(sHdr.Len)
	fmt.Println(*(*byte)(unsafe.Pointer(sHdr.Data)))

	// b := *(*[]byte)(unsafe.Pointer(sHdr))

	sl := &reflect.SliceHeader{
		Data: sHdr.Data,
		Len:  sHdr.Len,
		Cap:  sHdr.Len,
	}
	b := *(*[]byte)(unsafe.Pointer(sl))

	// bd := *(*unsafe.Pointer)(unsafe.Pointer(sHdr.Data))
	// fmt.Println(bd)
	// b := *(*[]byte)(bd)

	// fmt.Printf("%p\n", sHdr)
	fmt.Println(len(b), cap(b))
	fmt.Println(string(b[:len(s)]))
	fmt.Println(b[:len(s)])

	bs := make([]struct{}, int(1<<46)) // allocate struct{} which takes no memory
	// bs := make([]int64, int(1<<46)) // panic: runtime error: makeslice: len out of range
	fmt.Println(len(bs))
	// bs[len(bs)-1] = 1
	// fmt.Println(len(bs))

	var x byte
	fmt.Println(unsafe.Sizeof(x))
	var y struct{}
	fmt.Println(unsafe.Sizeof(y))
	fmt.Printf("%x\n", maxUintptr)
}

const maxUintptr = ^uintptr(0)
