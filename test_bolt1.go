package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"os"
)

func main() {
	os.Remove("/tmp/test1.db")

	b, err := bolt.Open("/tmp/test1.db", 0600, nil)
	fmt.Println(b, err)
	fmt.Println("Writing data.")

	size := 124
	maxMmapStep := 10
	sz := int64(size)
	if remainder := sz % int64(maxMmapStep); remainder > 0 {
		sz += int64(maxMmapStep) - remainder
	}
	fmt.Println(sz)
}
