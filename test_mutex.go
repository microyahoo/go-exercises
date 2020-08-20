package main

import (
	"fmt"
	"sync"
)

var _ sync.WaitGroup

type mutex struct {
	a, b int
}

func (m *mutex) printx() string {
	return "mutex"
}

func main() {
	var m *mutex
	fmt.Println(m)
	fmt.Println(m.printx())
}
