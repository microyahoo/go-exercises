package main

import (
	"fmt"
)

type meta struct {
	magic    uint32
	version  uint32
	pageSize uint32
}

func (m *meta) copy(dst *meta) {
	*dst = *m
}

func (m *meta) String() string {
	return fmt.Sprintf("meta[magic = %v, version = %v, pageSize = %v]", m.magic, m.version, m.pageSize)
}

func main() {
	m1 := &meta{
		magic:    12,
		version:  2,
		pageSize: 4096,
	}
	m2 := &meta{}
	m1.copy(m2)
	fmt.Println(m1)
	fmt.Println(m2)
	m2.magic = 18
	fmt.Println(m1)
	fmt.Println(m2)
}
