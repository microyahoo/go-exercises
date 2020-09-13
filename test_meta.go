package main

import (
	"fmt"
	"hash/fnv"
	"unsafe"
)

type txid uint64

type meta struct {
	magic    uint32
	version  uint32
	pageSize uint32
	checksum uint64
}

// generates the checksum for the meta.
func (m *meta) sum64() uint64 {
	var h = fnv.New64a()
	_, _ = h.Write((*[unsafe.Offsetof(meta{}.checksum)]byte)(unsafe.Pointer(m))[:])
	return h.Sum64()
}

func (m *meta) copy(dst *meta) {
	*dst = *m
}

func (m *meta) String() string {
	return fmt.Sprintf("meta[magic = %v, version = %v, pageSize = %v, checksum = %v]", m.magic, m.version, m.pageSize, m.checksum)
}

func main() {
	m1 := &meta{
		magic:    12,
		version:  2,
		pageSize: 4096,
	}
	fmt.Println(unsafe.Offsetof(meta{}.checksum))
	fmt.Println((*[unsafe.Offsetof(meta{}.checksum)]byte)(unsafe.Pointer(m1))[:])
	m1.checksum = m1.sum64()
	m2 := &meta{}
	m1.copy(m2)
	fmt.Println(m1)
	fmt.Println(m2)

	m2.magic = 18
	fmt.Println(m1)
	fmt.Println(m2)

	var x = 12
	var y = 24
	fmt.Printf("%X\n", x^y)

	minid := txid(0xFFFFFFFFFFFFFFFF)
	fmt.Println(minid)
	fmt.Println(minid > 0)

	// var m = 2.035013634110407
	// var n = 4.23188836343827
	// fmt.Printf("%X\n", m^n)
}
