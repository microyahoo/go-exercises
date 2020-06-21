package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"strings"
)

func main() {
	var x uint64
	x = 0x01204030
	buf := make([]byte, 8)
	n := binary.PutUvarint(buf, x)
	fmt.Println(x, n, buf)

	var y uint32
	y = 0x01203040
	buf2 := make([]byte, 0)
	buf2 = append(buf2, byte(y>>24), byte(y>>16), byte(y>>8), byte(y))
	fmt.Println(buf2)

	binary.LittleEndian.PutUint64(buf, x)
	fmt.Println(x, buf)

	binary.BigEndian.PutUint64(buf, x)
	fmt.Println(x, buf)

	s := "00000012"
	var num int64
	// ascii code
	fmt.Println([]byte(s))
	reader := strings.NewReader(s)
	err := binary.Read(reader, binary.LittleEndian, &num)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(num)

	s = string([]byte{0, 0, 0, 0, 0, 0, 0, 12})
	fmt.Println([]byte(s))
	// reader = strings.NewReader(s)
	reader.Reset(s)
	err = binary.Read(reader, binary.LittleEndian, &num)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(num)
}
