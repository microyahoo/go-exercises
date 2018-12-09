package main

import "fmt"

func main() {
	var ch = make(chan int)
	close(ch) //UNCOMMENT IT

	var ch2 = make(chan int)
	go func() {
		for i := 1; i < 10; i++ {
			ch2 <- i
		}
		close(ch2)
	}()

	for i := 0; i < 20; i++ {
		select {
		case x, ok := <-ch:
			fmt.Println("closed", x, ok)
		case x, ok := <-ch2:
			fmt.Println("open", x, ok)
		}
	}
}
