package linkname

import _ "unsafe"

//go:linkname hello github.com/microyahoo/go-exercises/runtime/outer.World
func hello() {
	println("hello, world!")
}
