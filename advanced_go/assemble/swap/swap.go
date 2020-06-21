package main

import (
	"fmt"
	"runtime"
)

func Swap(int, int) (int, int)

func test()
func printnl_nosplit()

func println(t int) {
	fmt.Println(t)
}

func main() {
	fmt.Println(Swap(1, 2))

	var local [1]struct {
		a bool
		b int16
		c []byte
	}

	var sp = &local[0]
	fmt.Println(sp)

	fmt.Println("------------")
	test()
	printnl_nosplit()
	fmt.Println("------------")

	for skip := 0; ; skip++ {
		pc, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		p := runtime.FuncForPC(pc)
		fnfile, fnline := p.FileLine(0)

		fmt.Printf("skip = %d, pc = 0x%08X\n", skip, pc)
		fmt.Printf("   func: file = %s, line = L%03d, name = %s, entry = 0x%08X\n", fnfile, fnline, p.Name(), p.Entry())
		fmt.Printf("   call: file = %s, line = L%03d\n\n", file, line)
	}
}
