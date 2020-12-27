package main

import (
	"fmt"
	"runtime"
)

func showData(txt string) {
	for i := 0; i < 5; i++ {
		runtime.Gosched()
		fmt.Println("RealOutput==>", txt)
	}
}

func main() {
	fmt.Println("stage 1")
	go showData("hi")
	fmt.Println("stage 2")
	showData("all")
	fmt.Println("stage 3")
}
