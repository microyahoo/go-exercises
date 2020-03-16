package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3*time.Second))
	defer cancel()

	go func(ctx context.Context) {
		defer cancel()

		// simulate a process that takes 2 second to complete
		time.Sleep(2 * time.Second)
	}(ctx)

	select {
	case <-ctx.Done():
		switch ctx.Err() {
		case context.DeadlineExceeded:
			fmt.Println("context timeout exceeded")
		case context.Canceled:
			fmt.Println("cancel the context by force")
		}
	}
}
