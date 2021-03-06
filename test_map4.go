package main

import (
	"fmt"
	"unsafe"
)

type programmer struct {
	name     string
	language string
}

func getMap() map[string]int {
	return nil
}

func main() {
	var mp map[string]int
	fmt.Println(mp == nil)
	fmt.Println(mp["x"])
	fmt.Println(getMap() == nil)
	x, ok := getMap()["x"]
	fmt.Println(x, ok)
	fmt.Println(getMap()["x"])
	mp = make(map[string]int)
	mp["qcrao"] = 100
	mp["stefno"] = 18

	count := **(**int)(unsafe.Pointer(&mp))
	fmt.Println(count, len(mp)) // 2 2

	p := programmer{"stefno", "go"}
	fmt.Println(p)

	name := (*string)(unsafe.Pointer(&p))
	*name = "qcrao"

	lang := (*string)(unsafe.Pointer(uintptr(unsafe.Pointer(&p)) + unsafe.Offsetof(p.language)))
	*lang = "Golang"

	fmt.Println(p)

	fmt.Println("-----------------------------------")
	rules := map[string]int{
		"a": 2,
		"b": 3,
	}
	fmt.Println(rules)
	for k, v := range rules {
		v++
		rules[k] = v
	}
	fmt.Println(rules)

	var programmers []*programmer
	fmt.Println(programmers)
	programmers = programmers[:0]
	fmt.Println(programmers)

	fmt.Println("-----------------------------------")
	for k := range rules {
		fmt.Println(k)
	}
}
