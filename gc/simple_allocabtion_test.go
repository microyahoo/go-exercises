package main

import (
	// "fmt"
	"runtime"
	"testing"
	"time"
)

// func main() {
// 	fmt.Println("vim-go")
// }

// BenchmarkAllocationEveryMs ...
func BenchmarkAllocationEveryMs(b *testing.B) {
	// need permanent allocation to clear see when the heap double its size
	var s *[]int
	tmp := make([]int, 1100000, 1100000)
	s = &tmp

	var a *[]int
	for i := 0; i < b.N; i++ {
		tmp := make([]int, 10000, 10000)
		a = &tmp

		time.Sleep(time.Millisecond)
	}
	_ = a
	runtime.KeepAlive(s)
}

// https://medium.com/a-journey-with-go/go-how-does-the-garbage-collector-watch-your-application-dbef99be2c35
