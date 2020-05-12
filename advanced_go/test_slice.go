package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func main() {
	var a []int
	var c = []int{}
	fmt.Println(a == nil)
	fmt.Println(c == nil)

	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&a))
	fmt.Printf("slice length: %d, capacity: %d, data: %v\n", sliceHeader.Len, sliceHeader.Cap, sliceHeader.Data)
	sliceHeader = (*reflect.SliceHeader)(unsafe.Pointer(&c))
	fmt.Printf("slice c length: %d, capacity: %d, data: %v\n", sliceHeader.Len, sliceHeader.Cap, sliceHeader.Data)
	fmt.Println("**************************")
	for i := 0; i < 20; i++ {
		a = append(a, i)
		sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		fmt.Printf("slice length: %d, capacity: %d, data: %v\n", sliceHeader.Len, sliceHeader.Cap, sliceHeader.Data)
	}
	fmt.Println("**************************")
	sliceHeader = (*reflect.SliceHeader)(unsafe.Pointer(&c))
	fmt.Printf("slice c length: %d, capacity: %d, data: %v\n", sliceHeader.Len, sliceHeader.Cap, sliceHeader.Data)
	c = a
	fmt.Println("**************************")

	for len(a) > 0 {
		a = a[1:]
		sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&a))
		fmt.Printf("slice length: %d, capacity: %d, data: %v\n", sliceHeader.Len, sliceHeader.Cap, sliceHeader.Data)
	}
	fmt.Println("**************************")
	sliceHeader = (*reflect.SliceHeader)(unsafe.Pointer(&c))
	fmt.Printf("slice c length: %d, capacity: %d, data: %v\n", sliceHeader.Len, sliceHeader.Cap, sliceHeader.Data)
}
