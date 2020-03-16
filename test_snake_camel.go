package main

import (
	"fmt"
	"github.com/iancoleman/strcase"
)

func main() {
	x := "TestCaseHelloWorld"
	y := strcase.ToSnake(x)
	fmt.Println(y)
}
