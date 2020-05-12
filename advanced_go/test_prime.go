package main

import (
	"context"
	"fmt"
)

func generateNatural(ctx context.Context) chan int {
	ch := make(chan int)

	go func() {
		for i := 2; ; i++ {
			select {
			case <-ctx.Done():
				return
			case ch <- i:
			}
		}
	}()
	return ch
}

func primeFilter(ctx context.Context, in <-chan int, prime int) chan int {
	out := make(chan int)

	go func() {
		for {
			if i := <-in; i%prime != 0 {
				fmt.Println("\t**", i, prime, in)
				select {
				case <-ctx.Done():
					return
				case out <- i:
				}
			}
		}
	}()

	return out
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	ch := generateNatural(ctx)

	for i := 0; i < 10; i++ {
		prime := <-ch
		fmt.Printf("%v: %v\n", i+1, prime)
		ch = primeFilter(ctx, ch, prime)
	}
	cancel()
}
