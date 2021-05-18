package main

import (
	"fmt"
	"unsafe"
)

type MyStruct struct {
	field int
}

func (self MyStruct) modify_struct_value() {
	fmt.Println(unsafe.Pointer(&self))
	self.field = 2
}

func (self MyStruct) copy_my_self() MyStruct {
	fmt.Println(unsafe.Pointer(&self))
	return self
}

func main() {
	var m1 MyStruct
	var m2 *MyStruct
	fmt.Println(m1)
	fmt.Println(m2 == nil)
	m2 = &MyStruct{1}
	fmt.Println(m2 == nil)
	local_variable := MyStruct{1}
	fmt.Println(unsafe.Pointer(&local_variable))
	local_variable.modify_struct_value()
	fmt.Println(local_variable) // {1}
	copied := local_variable.copy_my_self()
	fmt.Println(unsafe.Pointer(&copied))
	fmt.Println("---------------")
	var x, y int
	// x = y = 1
	fmt.Println(x, y)
}

// https://segmentfault.com/a/1190000006803598
