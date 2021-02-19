package main

import (
	"fmt"
	"time"
)

func main() {
	stream := make(chan string)
	start := time.Now()
	go func() {
		time.Sleep(time.Second * 3)
		stream <- "hello world"
	}()
	salutation, ok := <-stream
	fmt.Printf("(%v): %v: %v\n", ok, salutation, time.Since(start))

	intStream := make(chan int)
	go func() {
		for i := 1; i < 5; i++ {
			intStream <- i
			time.Sleep(time.Second)
		}
	}()

	for i := range intStream {
		fmt.Println(i)
	}
}
