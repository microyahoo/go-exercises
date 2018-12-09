package main

import (
	"fmt"
	"time"
)

func worker(done chan<- struct{}) {
	time.Sleep(time.Second * 5)
	// done <- struct{}{}
	close(done)
}

func main() {
	done := make(chan struct{})
	tick := time.NewTicker(time.Second * 1)

	go worker(done)

	for {
		select {
		case <-tick.C:
			fmt.Printf(".")
		case <-done:
			fmt.Println("done.")
			return
			// done = nil
		}
	}
}
