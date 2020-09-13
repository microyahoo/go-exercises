package main

// import (
// 	"fmt"
// 	"hash/fnv"
// 	"unsafe"
// )

// func main() {
// 	fmt.Println("vim-go")

// }

// type meta struct {
// 	magic    uint32
// 	version  uint32
// 	pageSize uint32
// 	flags    uint32
// 	checksum uint64
// }

// // generates the checksum for the meta.
// func (m *meta) sum64() uint64 {
// 	var h = fnv.New64a()
// 	_, _ = h.Write((*[unsafe.Offsetof(meta{}.checksum)]byte)(unsafe.Pointer(m))[:])
// 	return h.Sum64()
// }
