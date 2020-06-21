package main

import "fmt"

func main() {
	fmt.Println("vim-go")
	var ref Interface = &Impl{}
	ref.Do()
}

type Interface interface {
	Do()
}

type Impl struct{}

func (this *Impl) Do() {
}
