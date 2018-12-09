package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(name string) {
	for i := 0; i < 10; i++ {
		fmt.Println(name)
		time.Sleep(time.Millisecond * 100)
	}
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		worker("hello-------")
		wg.Done()
	}()

	go func() {
		worker("world***********")
		wg.Done()
	}()

	wg.Wait()
}
