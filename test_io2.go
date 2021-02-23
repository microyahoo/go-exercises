package main

import (
	"bufio"
	"fmt"
	"io"
	"time"
)

func main() {
	r, w := io.Pipe()

	fmt.Println(int(^uint(0) >> 1))
	go func() {
		writer := bufio.NewWriter(w)
		for i := 0; i < 3; i++ {
			time.Sleep(time.Second)
			writer.WriteString(time.Now().String() + "\n")
			writer.Flush()
		}
		// w.Close()
	}()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	// r.Close()
}
