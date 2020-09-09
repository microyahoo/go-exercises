package main

import (
	"fmt"
	"math"
	"os"
	"runtime"
)

func main() {
	fmt.Println(runtime.GOOS)
	fmt.Println(os.Getpagesize())
	fmt.Println(0x1000)
	fmt.Println(math.Pow(2, 12))
}
