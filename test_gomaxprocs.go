package main

import (
	"fmt"
	"runtime"
)

func init() {
	runtime.GOMAXPROCS(4) //使用多核
}

func main() {
	exit := make(chan int)
	go func() {
		defer close(exit)
		go func() {
			fmt.Println("b")
		}()
	}()

	for i := 0; i < 4; i++ {
		fmt.Println("a:", i)

		if i == 1 {
			runtime.Gosched() //切换任务
		}
	}
	<-exit
}

// https://www.jianshu.com/p/e45cea3e1723
