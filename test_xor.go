package main

import (
	"fmt"
	"math"
)

func main() {
	v1 := math.Float64bits(117)
	v2 := math.Float64bits(118)
	fmt.Printf("%064b\n", v1)
	fmt.Printf("%064b\n", v2)
	fmt.Println("-----------")
	fmt.Printf("%064b\n", v1^v2)
}
