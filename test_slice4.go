package main

import (
	"fmt"
)

func main() {
	var str1 []string
	var str2 []string

	str1 = append(str1, "a", "b")
	str2 = str1
	fmt.Println(str1, len(str1), cap(str1))
	fmt.Println(str2, len(str2), cap(str2))
	str1 = append(str1, "c", "d")
	fmt.Println(str1, len(str1), cap(str1))
	fmt.Println(str2, len(str2), cap(str2))
}
