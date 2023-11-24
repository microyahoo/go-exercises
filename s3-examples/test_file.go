package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	f, err := os.Open("test-upload")
	if err != nil {
		panic(err)
	}
	bytes, err := io.ReadAll(f)
	if err != nil {
		panic(err)
	}
	fmt.Println(len(bytes))

	for i := 0; i < len(bytes); {
		fmt.Printf("%s    ", string(bytes[i]))
		i += 1024 * 1024
	}
	fmt.Println()
}
