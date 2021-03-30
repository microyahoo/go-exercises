package main

import (
	"fmt"
)

func main() {
	var s1 []int
	s2 := []int{1, 2, 3}
	s3 := []int{4, 5, 6, 7}
	s4 := []int{1, 2, 3}
	fmt.Println(s2[:])
	fmt.Println(s2[:len(s2)-1])
	fmt.Println(s3[:])
	// 1、
	n1 := copy(s1, s2)
	fmt.Printf("n1=%d, s1=%v, s2=%v\n", n1, s1, s2)
	fmt.Println("s1 == nil", s1 == nil)
	// 2、
	n2 := copy(s2, s3)
	fmt.Printf("n2=%d, s2=%v, s3=%v\n", n2, s2, s3)
	// 3、
	n3 := copy(s3, s4)
	fmt.Printf("n3=%d, s3=%v, s4=%v\n", n3, s3, s4)

	fmt.Println(BlockMigrationStatusFail)

	// 切片在函数间以值的方式传递。由于切片的尺寸很小（在 64 位架构的机器上，一个切片需要 24 字节的内存：指针字段、长度和容量字段各需要 8 字节），在函数间复制和传递切片成本也很低。切片发生复制时，底层数组不会被复制，数组大小也不会有影响。

	a := []int{0, 1, 2, 3}
	FuncSlice(a, 4)
	fmt.Println(a)

	const XX = 1 << iota
	fmt.Println(XX)
	fmt.Println(^uintptr(0))
	fmt.Println("-------------1")
	fmt.Println(^1)
	fmt.Println(^uint64(1))
	fmt.Println(^uint32(1))
	fmt.Println("-------------2")
	fmt.Println(^int64(1))
	fmt.Println(^int32(1))
	// fmt.Println(int(^uintptr(0)))
	test(nil)
}

func test(rules []*Rule) {
	if len(rules) == 0 {
		fmt.Println(0, "hello")
		return
	}
	fmt.Println(1)
}

// Rule ...
type Rule struct{}

const (
	BlockMigrationStatusFree = iota
	BlockMigrationStatusFail
)

func FuncSlice(s []int, t int) {
	s[0]++
	s = append(s, t)
	s[0]++
}
