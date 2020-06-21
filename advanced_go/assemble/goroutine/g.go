package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

func getg() unsafe.Pointer

func getGoroutineId() int64 {
	g := getg()
	var g_goid_offset uintptr
	if f, ok := reflect.TypeOf(g).FieldByName("goid"); ok {
		g_goid_offset = f.Offset
	} else {
		panic("Can not find g.goid field")
	}
	p := (*int64)(unsafe.Pointer(uintptr(g) + g_goid_offset))
	return *p
}

func main() {
	fmt.Println("vim-go")
	fmt.Println(getGoroutineId())
}
