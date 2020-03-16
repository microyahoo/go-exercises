package main

import (
	"fmt"
)

const GroupName = "andy"

const (
	ErrCodeRbd = iota + 2
	ErrCodeVolume
)

var X = 100

type hello struct{}

func test() {
	fmt.Println("hello")
	alerts := []int{1, 2}
	var x []hello
	for i := 1; i < len(alerts); i++ {
		fmt.Println(alerts[i])
	}
	fmt.Println(x == nil)
	fmt.Println(x)
	fmt.Println(len(x) == 0)
}

func main() {
	test()
}
