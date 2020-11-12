package main

// https://github.com/golang/go/issues/37942

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGURG)

	for {
		select {
		case sig := <-ch:
			fmt.Printf("received %v: %s\n", sig, time.Now())
		default:
			_ = new(int) // generate some GC activities
		}
	}
}
