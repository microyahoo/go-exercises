package main

import "fmt"

func binarySearch(a []int64, value int64) {
	lo := 0
	hi := len(a) - 1
	for lo <= hi {
		mid := lo + (hi-lo)/2
		if a[mid] < value {
			lo = mid + 1
		} else if a[mid] > value {
			hi = mid - 1
		} else {
			fmt.Println("Success", mid)
			return
		}
	}
	fmt.Println("Not found", value)
}

func main() {
	a := []int64{1, 2, 3, 4}
	binarySearch(a, 1)
	binarySearch(a, 3)
	binarySearch(a, 5)

	var ch *int
	var ready []*int
	ready = make([]*int, 10)
	var b = 1
	ready[1] = &b
	ch = ready[1]
	ready[1] = nil
	fmt.Printf("%v, %v, %v\n", ch, *ch, ready[1])
}
