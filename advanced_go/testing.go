package main

import (
	"fmt"
	"testing"
)

type tb struct {
	testing.TB
}

func (p *tb) Fatal(args ...interface{}) {
	fmt.Println("TB.Fatal disabled")
}

func main() {
	var t testing.TB = new(tb)
	t.Fatal("hello world")
}
