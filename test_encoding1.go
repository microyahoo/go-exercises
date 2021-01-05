package main

import (
	"encoding/binary"
	"fmt"
	"math"
)

func p(a int64) {
	b1 := make([]byte, 8)
	b2 := make([]byte, 8)

	binary.LittleEndian.PutUint64(b1, uint64(a))
	binary.BigEndian.PutUint64(b2, uint64(a))

	fmt.Printf("%0x\t%0x\n", b1, b2)
}

func main() {
	p(0)
	p(1)
	p(0xFFFF0000)
	p(math.MaxInt64)

	p(-1)
	p(-200)
	p(-300)
	p(math.MinInt64)

	fmt.Println("----------------------------------")

	d(0)

	d(1)
	d(2)

	d(-1)
	d(-2)

	fmt.Println("----------------------------------")
	x := float64ToUint64(-1)
	// fmt.Printf("%0x\n", x)
	fmt.Println(x)
	y := uint64ToFloat64(x)
	// fmt.Printf("%0x\n", y)
	fmt.Println(y)

	x = float64ToUint64(1)
	fmt.Println(x)
	y = uint64ToFloat64(x)
	fmt.Println(y)
}

func d(a float64) {
	b := make([]byte, 8)

	bits := math.Float64bits(a)
	binary.BigEndian.PutUint64(b, bits)

	fmt.Printf("%0x\n", b)
	fmt.Printf("%0x\n", ^bits)
	fmt.Printf("%064b\n", binary.BigEndian.Uint64(b))
}

// We can not use lexicographically bytes comparison for negative and positive float directly.
// so here we will do a trick below.
func float64ToUint64(f float64) uint64 {
	u := math.Float64bits(f)
	if f >= 0 {
		u |= 0x8000000000000000
	} else {
		u = ^u
	}
	return u
}

func uint64ToFloat64(u uint64) float64 {
	fmt.Printf("%v, %t\n", u&0x8000000000000000, u&0x8000000000000000 > 0)
	if u&0x8000000000000000 > 0 {
		u &= ^uint64(0x8000000000000000)
	} else {
		u = ^u
	}
	return math.Float64frombits(u)
}

// https://www.jianshu.com/p/edb0a016e477

// Output:
// 0000000000000000        0000000000000000
// 0100000000000000        0000000000000001
// 0000ffff00000000        00000000ffff0000
// ffffffffffffff7f        7fffffffffffffff
// ffffffffffffffff        ffffffffffffffff
// 38ffffffffffffff        ffffffffffffff38
// d4feffffffffffff        fffffffffffffed4
// 0000000000000080        8000000000000000
// --------------------------------------------------------------------------
// 0000000000000000
// ffffffffffffffff
// 0000000000000000000000000000000000000000000000000000000000000000
// 3ff0000000000000
// c00fffffffffffff
// 0011111111110000000000000000000000000000000000000000000000000000
// 4000000000000000
// bfffffffffffffff
// 0100000000000000000000000000000000000000000000000000000000000000
// bff0000000000000
// 400fffffffffffff
// 1011111111110000000000000000000000000000000000000000000000000000
// c000000000000000
// 3fffffffffffffff
// 1100000000000000000000000000000000000000000000000000000000000000
