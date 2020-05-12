package main

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

func main() {
	defer debug.SetGCPercent(debug.SetGCPercent(-1))

	var count int32
	newfun := func() interface{} {
		return atomic.AddInt32(&count, 1)
	}

	pool := sync.Pool{New: newfun}

	v1 := pool.Get()
	fmt.Printf("v1 :%v\n", v1)

	pool.Put(9)
	pool.Put(10)
	pool.Put(11)
	pool.Put(12)

	v2 := pool.Get()
	fmt.Printf("v2 :%v\n", v2)

	debug.SetGCPercent(100)
	runtime.GC()

	v3 := pool.Get()
	fmt.Printf("v3 :%v\n", v3)

	pool.New = nil

	v4 := pool.Get()
	fmt.Printf("v4 :%v\n", v4)

	v5 := pool.Get()
	fmt.Printf("v5 :%v\n", v5)

	v6 := pool.Get()
	fmt.Printf("v6 :%v\n", v6)

	fmt.Println("\n\n")

	x := make([]int64, 8)
	fmt.Println(len(x))
	fmt.Println(cap(x))

	fmt.Println("\n\n")

	var w io.Writer
	w = os.Stdout
	var y nilTest
	fmt.Println(y.(nilTest))
}

type nilTest interface{}
