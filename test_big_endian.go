package main

import (
	"encoding/binary"
	"fmt"
)

func main() {
	var x uint64
	var c [binary.MaxVarintLen64]byte
	x = 123456
	fmt.Printf("%X\n", x)
	fmt.Println(c)
	binary.BigEndian.PutUint64(c[:], x)
	fmt.Println(c)
	// 1E240
	// [0 0 0 0 0 0 0 0 0 0]
	// [0 0 0 0 0 1 226 64 0 0]

	x = 1732
	fmt.Printf("type(x) = %T, x = %b\n", x, x)
	n := binary.PutUvarint(c[:], x)
	fmt.Println(c[:n])

	var y int64
	y = 1732
	var xxx int64
	xxx = -65
	fmt.Printf("type(y) = %T, y = %b, -65 = %b, uint64(-65) = %v, uint64(-65) = %b\n", y, y, xxx, uint64(xxx), uint64(xxx))
	n = binary.PutVarint(c[:], y)
	fmt.Println(c[:n])

	buf := make([]byte, binary.MaxVarintLen64)

	for _, x := range []int64{-65, 1, 2, 127, 128, 255, 256} {
		n := binary.PutVarint(buf, x)
		fmt.Print(x, " 输出的可变长度为：", n, "，十六进制为：")
		fmt.Printf("%v\n", buf[:n])
	}

	var t int64
	t = -65
	ux := uint64(t) << 1
	fmt.Printf("ux = %v, ux = %b, type(ux) = %T\n", ux, ux, ux)
	if t < 0 {
		ux = ^ux
	}
	fmt.Printf("ux = %v, ux = %b\n", ux, ux)

	fmt.Println("+++++++++++")
	bs := []byte{129, 1}
	fmt.Println(readUvarintFromBytes(bs))
	ux, _ = readUvarintFromBytes(bs) // ok to continue in presence of error
	tx := int64(ux >> 1)
	fmt.Println(tx)
	fmt.Println(tx & 1)
	fmt.Println(ux & 1)
	// TODO why?
	if ux&1 != 0 {
		tx = ^tx
	}
	fmt.Println(tx)
	fmt.Printf("tx = %b, ^4=%v, ^15=%v\n", tx, ^4, ^15)

	var a byte = 0x0F
	fmt.Printf("%08b\n", a)
	fmt.Printf("%08b\n", ^a)
	fmt.Println(^a)

	test := 64
	fmt.Printf("type(test) = %T, test=%b, ^test=%b, ^test=%v\n", test, test, ^test, ^test)

	const xx = iota + 1
	fmt.Printf("xx = %v, type(xx) = %T\n", xx, xx)
}

func readUvarintFromBytes(b []byte) (uint64, error) {
	var x uint64
	var s uint
	for i := 0; i < len(b); i++ {
		if b[i] < 0x80 {
			if i > 9 || i == 9 && b[i] > 1 {
				return x, fmt.Errorf("overflow")
			}
			return x | uint64(b[i])<<s, nil
		}
		x |= uint64(b[i]&0x7f) << s
		s += 7
	}
	fmt.Println("not print")
	return x, nil
}

// https://blog.csdn.net/skh2015java/article/details/78451033
