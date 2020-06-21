package main

import "fmt"

var Id int64

func main() {
	fmt.Println(Id)
	idPtr := &Id
	*idPtr = 243
	fmt.Println(Id)
}
