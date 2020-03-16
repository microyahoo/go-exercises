package main

import (
	"fmt"
	"runtime"
)

func main() {
	call()
}

func a() string {
	return b()
}

func b() string {
	return c()
}
func c() string {
	for calldepth := 0; calldepth < 4; calldepth++ {
		fmt.Println(runtime.Caller(calldepth))
	}
	return "xyz"
}

func call() {
	fmt.Println(a())
}

// 17379235 /Users/xsky/go/src/github.com/microyahoo/go-exercises/test_runtime2.go 21 true
// 17379667 /Users/xsky/go/src/github.com/microyahoo/go-exercises/test_runtime2.go 17 true
// 17379662 /Users/xsky/go/src/github.com/microyahoo/go-exercises/test_runtime2.go 13 true
// 17379661 /Users/xsky/go/src/github.com/microyahoo/go-exercises/test_runtime2.go 27 true
// xyz
