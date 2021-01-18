package main

import (
	"fmt"
	"time"
)

func main() {
	var (
		chann chan []byte
	)
	chann = make(chan []byte)

	go func(ch chan []byte) {
		x := <-ch
		fmt.Println(string(x))
	}(chann)

	chann <- []byte("xxyyzzafddddddddddddddddddmmnnqqoo")

	time.Sleep(2 * time.Second)

	done := make(chan error)
	waitNotifyC := make(chan struct{})

	go func() {
		select {
		case done <- doSomething():
			fmt.Printf("Receive do something\n")
		case <-waitNotifyC:
			fmt.Printf("Receive wait notify\n")
		}
	}()

	// ticker := time.NewTicker(time.Second)
	// defer ticker.Stop()
	select {
	// case <-ticker.C:
	case <-time.After(5 * time.Second):
		close(waitNotifyC)
		fmt.Printf("close wait notify\n")
	}
	time.Sleep(time.Second * 2)
	fmt.Println("Done...")
}

func doSomething() error {
	time.Sleep(30 * time.Second)
	return nil
}
