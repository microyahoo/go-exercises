package main

import (
	"fmt"
	"os"
)

func main() {

	c := &cmd{}
	for _, setupFd := range []F{(*cmd).stdin, (*cmd).stdout, (*cmd).stderr} {
		setupFd(c)
	}
}

type cmd struct {
}

func (c *cmd) stdout() (f *os.File, err error) {
	fmt.Println("stdout")
	return f, err
}

func (c *cmd) stdin() (f *os.File, err error) {
	fmt.Println("stdin")
	return f, err
}

func (c *cmd) stderr() (f *os.File, err error) {
	fmt.Println("stderr")
	return f, err
}

// F defines ...
type F func(*cmd) (*os.File, error)
