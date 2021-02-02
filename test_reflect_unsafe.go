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
	fmt.Println("-----------1------------")
	fmt.Println(unsafe.Sizeof(struct {
		a bool
		b int16
		c []int
		d func()
		e float64
	}{}))
	fmt.Println(unsafe.Sizeof(func() {}))
	fmt.Println(unsafe.Sizeof([]int{}))
	fmt.Println(unsafe.Sizeof([]byte{}))
	// fmt.Println(unsafe.Sizeof(io.Writer))
	var xi interface{}
	y := struct{}{}
	fmt.Println(unsafe.Sizeof(xi))
	fmt.Println(unsafe.Sizeof(y))

	fmt.Println("-----------2------------")
	u := new(user)
	fmt.Println(*u)

	pName := (*string)(unsafe.Pointer(u))
	*pName = "张三"

	pAge := (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(u)) + unsafe.Offsetof(u.age)))
	*pAge = 20

	fmt.Println(*u)
	fmt.Println(0xc4200c3010)

	fmt.Println("-----------3------------")
	n := 3
	b := make([]byte, n)
	b = []byte("abc")
	end := (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(&b[0])) + uintptr(n-1)))
	fmt.Println(*end)
	fmt.Println(b)
}

type user struct {
	name string
	age  int
}
