package main

import (
	"bytes"
	"fmt"
	"reflect"
	"unsafe"
)

func myAppend(s []int) []int {
	// 这里 s 虽然改变了，但并不会影响外层函数的 s
	s = append(s, 100)
	return s
}

func myAppendPtr(s *[]int) {
	// 会改变外层 s 本身
	*s = append(*s, 100)
	return
}

type path []byte

func (p path) TruncateAtFinalSlash() {
	i := bytes.LastIndex(p, []byte("/"))
	if i >= 0 {
		p = p[0:i]
	}
}

func main() {
	pathName := path("/usr/bin/tso") // Conversion from string to path.
	pathName.TruncateAtFinalSlash()
	fmt.Printf("%s\n", pathName)

	// s := make([]int, 4)
	s := make([]int, 0, 4)
	fmt.Println(len(s))
	fmt.Println(cap(s))
	newS := myAppend(s)

	fmt.Println("s=", s)
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer)(&s))
	fmt.Println(sliceHeader.Cap)
	fmt.Println(sliceHeader.Len)
	fmt.Println(sliceHeader.Data)
	sliceHeader2 := new(reflect.SliceHeader)
	sliceHeader2.Data = sliceHeader.Data
	sliceHeader2.Len = sliceHeader.Cap
	sliceHeader2.Cap = sliceHeader.Cap
	slice2 := (*[]int)(unsafe.Pointer(sliceHeader2))
	fmt.Println(*slice2)
	fmt.Println("newS=", newS)

	s = newS

	myAppendPtr(&s)
	fmt.Println(s)

	fmt.Println("------------------------------------------------------------")
	x := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	fmt.Printf("x = %v\n", x)
	fmt.Printf("x[8:] = %v\n", x[8:])
	fmt.Printf("x[5:8] = %v\n", x[5:8])
}
