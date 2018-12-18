package main

import (
	"fmt"
	"os"
	"sort"
	"unsafe"
)

type ints []int

func (a ints) merge(b ints) ints {
	if len(a) == 0 {
		return b
	}
	if len(b) == 0 {
		return a
	}
	merged := make(ints, len(a)+len(b))
	mergeints(merged, a, b)
	return merged
}

func mergeints(dst, a, b ints) {
	// Slicing does not copy the slice's data. It creates a new slice value that points to the original array.
	// https://blog.golang.org/go-slices-usage-and-internals
	merged := dst[:0]
	// merged := make(ints, 0)
	fmt.Printf("The merged is %v, &merged is %v\n", merged, uintptr(unsafe.Pointer(&merged)))
	fmt.Printf("The dst is %v, &dst is %v\n", dst, uintptr(unsafe.Pointer(&dst)))

	lead, follow := a, b
	if b[0] < a[0] {
		lead, follow = b, a
	}

	for len(lead) > 0 {
		index := sort.Search(len(lead), func(i int) bool {
			return lead[i] >= follow[0]
		})
		merged = append(merged, lead[:index]...)
		if index >= len(lead) {
			break
		}
		lead, follow = follow, lead[index:]
	}
	merged = append(merged, follow...)
}

func (s ints) Len() int           { return len(s) }
func (s ints) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ints) Less(i, j int) bool { return s[i] < s[j] }

func main() {
	a := ints{4, 8, 10, 12}
	b := ints{1, 2, 5, 7, 9, 13}
	c := a.merge(b)
	fmt.Println(c)
	fmt.Println(os.Getpagesize())
	step := 1 << 30
	size := 1 << 40
	fmt.Println(size % step)
	fmt.Println(100 % step)
}
