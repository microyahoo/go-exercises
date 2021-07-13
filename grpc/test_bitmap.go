package main

import (
	"fmt"
	"sync/atomic"
)

const smallBlockSize = 1024

type BitMap struct {
	Value  []uint64
	Length int64
}

func NewBitMap(length int64) *BitMap {
	var b BitMap
	b.Length = length
	tSlice := make([]uint64, length)
	b.Value = append(b.Value, tSlice...)

	_ := atomic.Value
	return &b
}

func (b *BitMap) expan(size int64) {
	if b.Length < smallBlockSize {
		tSlice := make([]uint64, b.Length)
		b.Value = append(b.Value, tSlice...)
	} else {
		tSlice := make([]uint64, b.Length/2)
		b.Value = append(b.Value, tSlice...)
	}
}

func (b *BitMap) Add(id int64) {
	floor, bit := id/64, id%64
	if floor > b.Length {
		b.expan(floor - b.Length)
	}
	b.Value[floor] = b.Value[floor] | uint64(1<<bit)
}

func (b *BitMap) Exist(id int64) bool {
	floor, bit := id/64, id%64
	if floor > b.Length {
		return false
	}
	return !(b.Value[floor]&uint64(1<<bit) == 0)
}

func main() {
	bp := NewBitMap(100)
	bp.Add(1)
	bp.Add(3)
	bp.Add(5)
	bp.Add(7)
	bp.Add(1000)

	fmt.Println(bp.Exist(1), bp.Exist(2), bp.Exist(3), bp.Exist(1000))
}
