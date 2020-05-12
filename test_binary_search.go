package main

import (
	"fmt"
)

func binarySearch(slice []int64, length int, val int64) (pos int) {
	if length == 0 || len(slice) == 0 {
		return -1
	}
	var i = length / 2
	fmt.Println(slice, i, length)
	if slice[i] > val {
		return binarySearch(slice[0:i], i, val)
	} else if slice[i] < val {
		return binarySearch(slice[i:], length-i, val)
	}
	return i
}

func binarySearch2(slice []int64, length int, val int64) (pos int) {
	if length == 0 || len(slice) == 0 {
		return -1
	}
	var low = 0
	var high = length - 1
	var mid int
	for {
		if low > high {
			return -1
		}
		mid = (low + high) / 2
		fmt.Println("++", slice, mid, length)
		if slice[mid] > val {
			high = mid - 1
		} else if slice[mid] < val {
			low = mid + 1
		} else {
			return mid
		}
	}
	return -1
}

func binarySearch3(slice []int64, length int, val int64) (pos int) {
	var low = 0
	var high = length - 1
	for low <= high {
		mid := low + (high-low)/2
		if val < slice[mid] {
			high = mid - 1
		} else if val > slice[mid] {
			low = mid + 1
		} else {
			return mid
		}
	}
	return -1
}

func main() {
	var slice = []int64{1, 3, 5, 7, 8, 9}
	fmt.Println(slice)
	fmt.Println(binarySearch(slice, 6, 8))
	fmt.Println(slice)
	fmt.Println("***", binarySearch2(slice, 6, 8))
	fmt.Println(slice)
	fmt.Println("***", binarySearch2(slice, 6, 88))
	fmt.Println("***", binarySearch2(slice, 6, 1))
	fmt.Println("***", binarySearch2(slice, 6, -8))
	fmt.Println("&+", binarySearch2(slice, 6, 88))
	fmt.Println("&+", binarySearch2(slice, 6, 1))
	fmt.Println("&+", binarySearch2(slice, 6, -8))
}
