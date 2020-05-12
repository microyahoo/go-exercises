package main

//void SayHello(_GoString_ s);
import "C"

import (
	"fmt"
)

func main() {
	C.SayHello("hello world")
}

// export SayHello
func SayHello(s string) {
	fmt.Println(s)
}
