package main

import (
	"context"
	"fmt"
	"log"
	"time"
)

func main() {
	// Context tree:
	// parent -> cancelChild
	// parent -> valueChild -> timerChild
	parent, cancel := context.WithCancel(context.Background())
	cancelChild, stop := context.WithCancel(parent)
	defer stop()

	valueChild := context.WithValue(parent, "key", "value")
	timerChild, stop := context.WithTimeout(valueChild, 10000*time.Hour)
	defer stop()
	cancelChild2, stop := context.WithCancel(timerChild)
	defer stop()

	cancel()

	// parent and children should all be finished.
	check := func(ctx context.Context, name string) {
		select {
		case <-ctx.Done():
		default:
			log.Fatalf("<-%s.Done() blocked, but shouldn't have", name)
		}
		if e := ctx.Err(); e != context.Canceled {
			log.Printf("%s.Err() == %v want %v", name, e, context.Canceled)
		}
	}
	check(parent, "parent")
	check(cancelChild, "cancelChild")
	check(valueChild, "valueChild")
	check(timerChild, "timerChild")
	check(cancelChild2, "cancelChild2")
	fmt.Println("hello")
}
