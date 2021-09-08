package main

import (
	"fmt"
	"math"
)

func main() {
	x := 0.210752
	fmt.Println(x * 100)
	fmt.Println(int64(x * 100))
	fmt.Println(math.Round(x * 100))
	fmt.Println("math.Ceil(99.0 / 8) = ", math.Ceil(99.0/8))
	fmt.Println("math.Round(99.0 / 8) = ", math.Round(99.0/8))
	fmt.Printf("%.2f\n", x*100)
	fmt.Printf("%b\n", 4<<10)

	var conns map[string]int
	if n := conns["a"]; n < 1 {
		fmt.Println(n, "xxx")
		conns["a"] = 1
	}
}
