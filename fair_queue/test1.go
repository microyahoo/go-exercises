package main

import (
	"fmt"
	"math"
)

const maxHashBits = 60
const spare = 64 - maxHashBits

// ff computes the falling factorial `n!/(n-m)!` and requires n to be
// positive and m to be in the range [0, n] and requires the answer to
// fit in an int
func ff(n, m int) int {
	ans := 1
	for f := n; f > n-m; f-- {
		ans *= f
	}
	return ans
}

func main() {
	fmt.Println(math.Log2(float64(ff(64, 3))))
	fmt.Println(math.Ceil(math.Log2(float64(ff(64, 3)))))
	fmt.Println(math.Ceil(math.Log2(float64(ff(64, 3)))) + spare)
	fmt.Println(1 << uint(math.Ceil(math.Log2(float64(ff(64, 3))))+spare))
	fmt.Println(ff(128, 3))
	fmt.Println(ff(50, 4))
	fmt.Println(ff(52, 4))
	fmt.Println(ff(54, 4))

	tests := []struct {
		deckSize, handSize int
		hashMax            int
	}{
		{64, 3, 1 << uint(math.Ceil(math.Log2(float64(ff(64, 3))))+spare)},
		{128, 3, ff(128, 3)},
		{50, 4, ff(50, 4)},
	}
	for _, test := range tests {
		fallingFactorial := ff(test.deckSize, test.handSize)
		permutations := ff(test.handSize, test.handSize)
		allCoordinateCount := fallingFactorial / permutations
		nff := float64(test.hashMax) / float64(fallingFactorial)
		minCount := permutations * int(math.Floor(nff))
		maxCount := permutations * int(math.Ceil(nff))
		fmt.Printf("fallingFactorial = %d, permutations = %d, allCoordinateCount = %d, nff = %f\n",
			fallingFactorial, permutations, allCoordinateCount, nff)
		fmt.Printf("minCount = %d\n", minCount)
		fmt.Printf("maxCount = %d\n", maxCount)
	}
}

// 17.93147623388679
// 18
// 22
// 4194304
// 2048256
// 5527200
// 6497400
// 7590024
// fallingFactorial = 249984, permutations = 6, allCoordinateCount = 41664, nff = 16.778290
// minCount = 96
// maxCount = 102
// fallingFactorial = 2048256, permutations = 6, allCoordinateCount = 341376, nff = 1.000000
// minCount = 6
// maxCount = 6
// fallingFactorial = 5527200, permutations = 24, allCoordinateCount = 230300, nff = 1.000000
// minCount = 24
// maxCount = 24
