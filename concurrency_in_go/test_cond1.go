package main

import (
	"fmt"
	"sync"
	"time"
)

// Button ...
type Button struct {
	Clicked *sync.Cond
}

func main() {
	button := Button{Clicked: sync.NewCond(&sync.Mutex{})}

	subscribe := func(c *sync.Cond, fn func()) {
		var goroutineRunning sync.WaitGroup
		goroutineRunning.Add(1)

		go func() {
			goroutineRunning.Done()
			c.L.Lock()
			defer c.L.Unlock()

			c.Wait()
			fn()
		}()

		goroutineRunning.Wait()
		fmt.Println("func done")
	}

	var clickRegistered sync.WaitGroup
	clickRegistered.Add(3)

	subscribe(button.Clicked, func() {
		fmt.Println("Maximizing window.")
		clickRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		fmt.Println("Displaying annoying dialog box!")
		clickRegistered.Done()
	})

	subscribe(button.Clicked, func() {
		fmt.Println("Mouse clicked.")
		clickRegistered.Done()
	})

	time.Sleep(time.Second)
	button.Clicked.Broadcast()

	clickRegistered.Wait()

	fmt.Println("--------------------------------------------")

	var onceA, onceB sync.Once
	var initB func()

	initA := func() {
		onceB.Do(initB)
	}
	initB = func() {
		onceA.Do(initA)
	}

	onceA.Do(initB)
}
