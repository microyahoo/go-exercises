package main

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func main() {
	_, filename, _, _ := runtime.Caller(1)
	fmt.Println(filename)
	fmt.Println(filepath.Join(filepath.Dir(filename), ".."))
	path, _ := filepath.Abs(filepath.Join(filepath.Dir(filename), ".."))
	fmt.Println(path)

	for skip := 0; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		fmt.Printf("skip = %v, pc = %v, file = %v, line = %v\n", skip, pc, file, line)
	}
	// Output:
	// skip = 0, pc = 17393793, file = /Users/xsky/go/src/github.com/microyahoo/go-exercises/test_runtime.go, line = 17
	// skip = 1, pc = 16947115, file = $(GOROOT)/src/runtime/proc.go, line = 200
	// skip = 2, pc = 17108944, file = $(GOROOT)/src/runtime/asm_amd64.s, line = 1337

	pc := make([]uintptr, 1024)
	for skip := 0; ; skip++ {
		n := runtime.Callers(skip, pc)
		if n <= 0 {
			break
		}
		fmt.Printf("skip = %v, pc = %v\n", skip, pc[:n])
	}
	// Output:
	// skip = 0, pc = [16806961 17394303 16947244 17109073]
	// skip = 1, pc = [17394303 16947244 17109073]
	// skip = 2, pc = [16947244 17109073]
	// skip = 3, pc = [17109073]
	oft := new(onlyForTest)
	fmt.Println(len(oft.es))
}

type onlyForTest struct {
	es []error
}
