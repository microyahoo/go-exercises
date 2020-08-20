package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

var mtx sync.Mutex
var changed bool

func main() {
	var timer1 = time.NewTicker(time.Second * 3)
	var timer2 = time.NewTicker(time.Second * 2)

	go func() {
		for {
			select {
			case <-timer2.C:
				fmt.Printf("***Change flag to %t\n", false)
				mtx.Lock()
				changed = false
				mtx.Unlock()
			default:
			}
		}
	}()

	go func() {
		for {
			select {
			case <-timer1.C:
				fmt.Printf("###Change flag to %t\n", true)
				mtx.Lock()
				changed = true
				mtx.Unlock()
			default:
			}
			if changed {
			} else {
				time.Sleep(3 * time.Second)
				keepPrimary()
			}
		}
	}()

	select {}
}

func keepPrimary() {
	monitorCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		fmt.Println("**monitoring metric service")
		for {
			select {
			case <-monitorCtx.Done():
				fmt.Println("exit monitor metric service")
				return
			default:
				fmt.Println("\tmonitor metric...")
				time.Sleep(5 * time.Second)
			}
		}
	}()
}
