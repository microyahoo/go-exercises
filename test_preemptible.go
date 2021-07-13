package main

import (
	"fmt"
	"runtime"
	"time"
)

func test() {
	a := 100
	for i := 1; i < 1000; i++ {
		a = i*100/i + a
	}
}

func main() {
	runtime.GOMAXPROCS(1)
	go func() {
		for {
			test()
		}
	}()
	time.Sleep(100 * time.Millisecond)
	fmt.Println("hello world")
}
