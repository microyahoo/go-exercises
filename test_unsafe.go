package main

import (
	"fmt"
	"unsafe"
)

type programmer struct {
	name     string
	language string
	test     *justTest
}

type justTest struct {
	a string
	b int
}

func main() {
	t := justTest{
		a: "hello",
		b: 100,
	}
	p := programmer{
		name:     "stefno",
		language: "go",
		test:     &t,
	}

	// {stefno go 0xc00007c020}
	// main.programmer{name:"stefno", language:"go", test:(*main.justTest)(0xc00007c020)}
	fmt.Println(p)
	fmt.Printf("%#v\n", p)

	name := (*string)(unsafe.Pointer(&p))
	*name = "qcrao"

	lang := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + unsafe.Offsetof(p.language)))
	*lang = "Golang"

	// {qcrao Golang 0xc00007c020}
	// 0xc000066180
	// main.programmer{name:"qcrao", language:"Golang", test:(*main.justTest)(0xc00007c020)}
	fmt.Println(p)
	fmt.Printf("%p\n", &p)
	fmt.Printf("%#v\n", p)
	// 将其转化为二级指针
	test := (*unsafe.Pointer)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + unsafe.Offsetof(p.test)))

	// 0xc0000661a0
	// &{ 17640398}
	// 0xc0000661a0
	// 0xc00007c020
	// &{hello 100}
	// test=*unsafe.Pointer, *test=unsafe.Pointer
	fmt.Println(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + unsafe.Offsetof(p.test)))
	fmt.Println((*justTest)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + unsafe.Offsetof(p.test)))) //invalid
	fmt.Println(test)
	fmt.Println(*test)
	fmt.Println((*justTest)(*test))
	fmt.Printf("test=%T, *test=%T\n", test, *test)
}
