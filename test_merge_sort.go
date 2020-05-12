package main

import (
	"fmt"
)

func main() {
	slice := []int{5, 1, 4, 8, 3, 9, 2, 7, 6}
	sortSlice(slice)
	fmt.Println(slice)
	data := []int{1, 2, 3}
	for _, v := range data {
		v *= 10 // data 中原有元素是不会被修改的
	}
	fmt.Println("data: ", data) // data:  [1 2 3]
	// https://studygolang.com/articles/12589

	var x struct {
		a bool
		b int16
		c []int
	}
	fmt.Printf("%p %p\n", &x, &x.a)
}

var aux []int

func sortSlice(a []int) {
	aux = make([]int, len(a))
	sort(a, 0, len(a)-1)
}

func sort(a []int, low, high int) {
	fmt.Printf("\tsort(a, %d, %d)\n", low, high)
	if low >= high {
		return
	}
	mid := low + (high-low)/2
	sort(a, low, mid)
	sort(a, mid+1, high)
	mergeSort(a, low, mid, high)
}

func mergeSort(a []int, low, mid, high int) {
	i := low
	j := mid + 1
	for k := low; k <= high; k++ {
		aux[k] = a[k]
	}
	for k := low; k <= high; k++ {
		if i > mid {
			a[k], j = aux[j], j+1
		} else if j > high {
			a[k], i = aux[i], i+1
		} else if aux[i] < aux[j] {
			a[k], i = aux[i], i+1
		} else {
			a[k], j = aux[j], j+1
		}
	}
}
