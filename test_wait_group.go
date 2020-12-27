package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for range ticker.C {
			fmt.Println("hello")
		}
	}()

	fmt.Println("xxx")
	time.Sleep(100 * time.Second)
}
