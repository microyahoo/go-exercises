package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)

	intStream := make(chan int, 4)
	go func() {
		defer close(intStream)
		// defer fmt.Println("Producer Done.")
		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
		for i := 0; i < 5; i++ {
			// fmt.Printf("Sending: %d\n", i)
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intStream <- i
			fmt.Fprintf(&stdoutBuff, "\t%d\n", i)
			// fmt.Printf("\t%d\n", i)
		}
	}()

	for integer := range intStream {
		fmt.Println("hello")
		// fmt.Printf("Received %v.\n", integer)
		fmt.Fprintf(&stdoutBuff, "Received %v.\n", integer)
	}
}
