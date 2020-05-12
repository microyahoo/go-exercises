package main

import (
	"fmt"
	"path/filepath"
	"sync"
	"sync/atomic"
)

var total uint64

func worker(wg *sync.WaitGroup) {
	defer wg.Done()

	var i uint64
	for i = 0; i <= 100; i++ {
		atomic.AddUint64(&total, i)
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	worker(&wg)
	worker(&wg)

	wg.Wait()
	fmt.Println(total)

	path := "/dev/installer/h:/usr/bin"
	fmt.Println(filepath.SplitList(path))

	// sync.Once
}
