package main

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 10)
	go func() {
		in := 1
		for {
			in++
			ch <- in
		}
	}()

	for {
		select {
		case _ = <-ch:
			// do something...
			fmt.Println("1111")
			continue
		case <-time_chan():
			// case <-time.After(3 * time.Minute): // 内存泄露
			fmt.Printf("现在是：%d！", time.Now().Unix())
		}
	}
}

func time_chan() <-chan time.Time {
	fmt.Println("hello world")
	return time.After(3 * time.Minute) // 内存泄露, 每次都会新建 Timer
}
