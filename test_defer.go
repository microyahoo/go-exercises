package main

import (
	"fmt"
)

func main() {
	a := "xxx"
	defer func() {
		fmt.Println(a)
	}()

	a = "yyy"
	fmt.Println("hello")
}
