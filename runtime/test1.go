package main

import (
	"time"

	"github.com/microyahoo/go-exercises/runtime/outer"
)

// https://www.pixelstech.net/article/1649596852-The-magic-of-go%3Alinkname
func main() {
	time.Sleep(time.Second)
	outer.World()
}
