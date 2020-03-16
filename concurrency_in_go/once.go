package main

import (
	"fmt"
	"sync"
)

func main() {
	var onceA, onceB sync.Once
	var initB func()

	initA := func() {
		fmt.Println("before **initA")
		onceB.Do(initB)
		fmt.Println("after **initA")
	}

	initB = func() {
		fmt.Println("before **initB")
		onceA.Do(initA)
		fmt.Println("after **initB")
	}

	onceA.Do(initA)
}
