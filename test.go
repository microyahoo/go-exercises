package main

import "fmt"

func binarySearch(a []int64, value int64) {
	lo := 0
	hi := len(a)
	for lo < hi {
		mid := lo + (hi-lo)/2
		if a[mid] < value {
			lo = mid
		} else if a[mid] > value {
			hi = mid
		} else {
			fmt.Println("Success")
			break
		}
	}
}

func main() {
	a := []int64{1, 2, 3, 4}
	binarySearch(a, 3)
}
