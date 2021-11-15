package main

import (
	"fmt"
	"sort"
)

type writeBuffer struct {
	keys []int
	used int
}

func (b *writeBuffer) Less(i, j int) bool {
	return b.keys[i] < b.keys[j]
}

func (b *writeBuffer) Len() int {
	return b.used
}

// Swap swaps the elements with indexes i and j.
func (b *writeBuffer) Swap(i, j int) {
	b.keys[i], b.keys[j] = b.keys[j], b.keys[i]
}

var _ sort.Interface = (*writeBuffer)(nil)

func main() {
	b := &writeBuffer{
		keys: []int{3, 6, 1, 7, 0, 4, 2, 22, 5, 17},
		used: 5,
	}
	fmt.Printf("%#v\n", b)
	sort.Sort(b)
	fmt.Printf("%#v\n", b)
}
