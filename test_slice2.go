package main

import (
	"bytes"
	"fmt"
)

func main() {
	var buffer [256]byte
	slice := buffer[10:20]
	for i := 0; i < len(slice); i++ {
		slice[i] = byte(i)
	}
	fmt.Println("before", slice)
	AddOneToEachElement(slice)
	fmt.Println("after", slice)

	fmt.Println("--------------------------")

	fmt.Println("Before: len(slice) =", len(slice))
	newSlice := SubtractOneFromLength(slice)
	fmt.Println("After:  len(slice) =", len(slice))
	fmt.Println("After:  len(newSlice) =", len(newSlice))

	fmt.Println("--------------------------")

	fmt.Println("Before: len(slice) =", len(slice))
	PtrSubtractOneFromLength(&slice)
	fmt.Println("After:  len(slice) =", len(slice))

	fmt.Println("--------------------------")

	pathName := path("/usr/bin/tso") // Conversion from string to path.
	pathName.PtrTruncateAtFinalSlash()
	fmt.Printf("%s\n", pathName)

	fmt.Println("--------------------------")

	pathName = path("/usr/bin/tso")
	pathName.ToUpper()
	fmt.Printf("%s\n", pathName)

	fmt.Println("--------------------------")
	// var iBuffer [10]int
	// iSlice := iBuffer[0:0]
	// for i := 0; i < 20; i++ {
	// 	iSlice = Extend(iSlice, i)
	// 	fmt.Println(iSlice)
	// }

	iSlice := make([]int, 0, 5)
	for i := 0; i < 10; i++ {
		iSlice = Extend(iSlice, i)
		fmt.Printf("len=%d cap=%d iSlice=%v\n", len(iSlice), cap(iSlice), iSlice)
		fmt.Println("address of 0th element:", &iSlice[0])
	}

	fmt.Println("--------------------------")

	pathName = path("/usr/bin/tso") // Conversion from string to path.
	pathName.TruncateAtFinalSlash()
	fmt.Printf("%s\n", pathName)

}

// AddOneToEachElement add one to each elements
func AddOneToEachElement(slice []byte) {
	for i := range slice {
		slice[i]++
	}
}

// SubtractOneFromLength ...
func SubtractOneFromLength(slice []byte) []byte {
	slice = slice[0 : len(slice)-1]
	return slice
}

// PtrSubtractOneFromLength ...
func PtrSubtractOneFromLength(slicePtr *[]byte) {
	slice := *slicePtr
	*slicePtr = slice[0 : len(slice)-1]
}

type path []byte

func (p *path) PtrTruncateAtFinalSlash() {
	i := bytes.LastIndex(*p, []byte("/"))
	if i >= 0 {
		*p = (*p)[0:i]
	}
}

func (p path) TruncateAtFinalSlash() {
	i := bytes.LastIndex(p, []byte("/"))
	if i >= 0 {
		p = p[0:i]
	}
}

func (p path) ToUpper() {
	for i, b := range p {
		if 'a' <= b && b <= 'z' {
			p[i] = b + 'A' - 'a'
		}
	}
}

// Extend ...
func Extend(slice []int, element int) []int {
	n := len(slice)
	if n == cap(slice) {
		// Slice is full; must grow.
		// We double its size and add 1, so if the size is zero we still grow.
		newSlice := make([]int, len(slice), 2*len(slice)+1)
		copy(newSlice, slice)
		slice = newSlice
	}
	slice = slice[0 : n+1]
	slice[n] = element
	return slice
}
