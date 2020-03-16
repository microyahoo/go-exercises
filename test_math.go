package main

import (
	"fmt"
	"math"
)

func main() {
	x := 0.210752
	fmt.Println(x * 100)
	fmt.Println(math.Round(x * 100))
	fmt.Printf("%.2f\n", x*100)
}
