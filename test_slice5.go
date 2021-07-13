package main

import "fmt"

func changeSlice(ints []int) {
	ints[0] = 10
}

func changeArray(ints [6]int) {
	fmt.Println(len(ints))
	fmt.Println(cap(ints))
	ints[0] = 10
	fmt.Println(ints[:1])
}

func main() {
	slice1 := []int{1, 2, 3, 4, 5, 6}
	array1 := [6]int{1, 2, 3, 4, 5, 6}
	fmt.Println(slice1)
	fmt.Println(array1)

	changeArray(array1)
	changeSlice(slice1)
	fmt.Println(slice1)
	fmt.Println(array1)

	slice := slice1[2:5]
	fmt.Println(slice)
	fmt.Println(len(slice))
	fmt.Println(cap(slice)) // The Capacity field is equal to the length of the underlying array, minus the index in the array of the first element of the slice (zero in this case).
}
