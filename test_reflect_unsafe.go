package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"unsafe"
)

func main() {
	fmt.Println(strings.Split("network__addresses", "__"))
	t := reflect.TypeOf(3)
	fmt.Println(t)
	fmt.Println(t.String())

	var w io.Writer = os.Stdout
	fmt.Println(reflect.TypeOf(w))

	v := reflect.ValueOf(3)
	fmt.Println(v)
	fmt.Printf("%v\n", v)

	fmt.Println(v.Kind())
	// fmt.Println(v.MapKeys())

	sl := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	slength := len(sl)
	xs := 0xFFFF
	fmt.Printf("the type of length is %T\n", slength)
	fmt.Printf("the type of 0xFFFF is %T, and the value is %v\n", xs, xs)
	fmt.Printf("sl[5:] = %v\n", sl[5:])
	fmt.Printf("sl[:] = %v\n", sl[:])
	fmt.Printf("sl[0:] = %v\n", sl[0:])

	fmt.Println(unsafe.Sizeof(string("xxx")))
	fmt.Println(unsafe.Sizeof(struct {
		bool
		float64
		int16
	}{}))
	fmt.Println(unsafe.Sizeof(struct {
		float64
		int16
		bool
	}{}))
	fmt.Println(unsafe.Sizeof(struct {
		bool
		int16
		float64
	}{}))
	fmt.Println(unsafe.Alignof(struct {
		bool
		int16
		float64
	}{}))

	var x struct {
		a bool
		b int16
		c []int
	}
	pb := (*int16)(unsafe.Pointer(uintptr(unsafe.Pointer(&x)) + unsafe.Offsetof(x.b)))
	*pb = 42
	fmt.Println(x.b)
	fmt.Println(0x7FFFFFF)
}
