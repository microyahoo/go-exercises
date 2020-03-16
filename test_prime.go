package main

import (
	"fmt"
)

func main() {
	ch := make(chan int)
	go genearte(ch)
	for i := 1; i < 10; i++ {
		prime := <-ch
		fmt.Println(prime)
		ch1 := make(chan int)
		go filter(ch, ch1, prime)
		ch = ch1
	}
}

func genearte(ch chan<- int) {
	for i := 2; ; i++ {
		fmt.Printf("genearte %d\n", i)
		ch <- i
	}
}

func filter(in <-chan int, out chan<- int, prime int) {
	for {
		i := <-in
		if i%prime != 0 {
			out <- i
		}
	}
}

// genearte 2
// genearte 3
// 2
// 3
// genearte 4
// genearte 5
// 5
// genearte 6
// genearte 7
// 7
// genearte 8
// genearte 9
// genearte 10
// genearte 11
// genearte 12
// genearte 13
// genearte 14
// genearte 15
// 11
// genearte 16
// genearte 17
