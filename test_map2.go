package main

import (
	"fmt"
)

func main() {
	a := make(map[int]string)
	a[1] = "a"
	a[2] = "b"
	b := get(a)
	fmt.Println(b)
	a[2] = "c"
	fmt.Println(b)
	for x := range b {
		fmt.Println(x)
	}
	fmt.Println(a[3])
}

func get(m map[int]string) map[int]string {
	return m
}
