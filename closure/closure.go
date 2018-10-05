package closure

import (
	"fmt"
)

func Squares() func() int {
	var x int
	fmt.Printf("**%d\n", x)
	return func() int {
		x++
		return x * x
	}
}
