package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-ticker.C:
			go panicFunc()
		}
	}
	select {}
}

func panicFunc() {
	time.Sleep(time.Second)
	panic(fmt.Sprintf("panic"))
}
