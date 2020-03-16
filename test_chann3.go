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
}
