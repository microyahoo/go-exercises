package main

import "fmt"
import "github.com/microyahoo/go-exercise/closure"

func main() {
	fmt.Println("vim-go")
	f := closure.Squares()
	fmt.Printf("%T\n", f)
	fmt.Println(f())
	fmt.Println(f())
	fmt.Println(f())
	fmt.Println(f())
}
