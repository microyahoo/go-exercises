package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	_ "path/filepath"
	"strings"
)

const (
	// PrivateFileMode ...
	PrivateFileMode = 0600
)

func main() {
	path := "/home/admin/data/node1/etcd/10.255.101.74.etcd/member/wal/0000000000000001-000000000006c314.wal"
	// path := "/Users/xsky/go/src/github.com/etcd-io/etcd/infra1.etcd/member/wal/0000000000000000-0000000000000000.wal"
	f, err := os.OpenFile(path, os.O_RDONLY, PrivateFileMode)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	reader := bufio.NewReader(f)
	var bytes = make([]byte, 8)
	var n int64
	x, err := reader.Read(bytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(x)
	fmt.Println(bytes)

	fmt.Println("===============1==============")
	_ = io.Copy
	f.Seek(0, io.SeekStart)
	reader = bufio.NewReader(f)
	err = binary.Read(reader, binary.LittleEndian, &n)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(n)
	fmt.Printf("%b\n", 0x84)
	fmt.Printf("%b\n", n)
	fmt.Printf("%X\n", n)

	recBytes, padBytes := decodeFrameSize(n)
	fmt.Println(recBytes, padBytes)

	fmt.Println("===============2==============")
	var num int64
	s := string([]byte{12, 0, 0, 0, 0, 0, 0, 0x80})
	fmt.Println([]byte(s))
	r := strings.NewReader(s)
	err = binary.Read(r, binary.LittleEndian, &num)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(num)
	fmt.Println(decodeFrameSize(num))
	fmt.Printf("%b\n", 0x84)

	fmt.Println("===============3==============")
	for i := 0; i < 20; i++ {
		encodeFrameSize(i)
	}
}

func decodeFrameSize(lenField int64) (recBytes int64, padBytes int64) {
	// the record size is stored in the lower 56 bits of the 64-bit length
	recBytes = int64(uint64(lenField) & ^(uint64(0xff) << 56))
	// non-zero padding is indicated by set MSb / a negative length
	if lenField < 0 {
		fmt.Println("lenField is negative")
		// padding is stored in lower 3 bits of length MSB
		padBytes = int64((uint64(lenField) >> 56) & 0x7)
	}
	return recBytes, padBytes
}

func encodeFrameSize(dataBytes int) (lenField uint64, padBytes int) {
	lenField = uint64(dataBytes)
	// force 8 byte alignment so length never gets a torn write
	padBytes = (8 - (dataBytes % 8)) % 8
	if padBytes != 0 {
		lenField |= uint64(0x80|padBytes) << 56
	}
	fmt.Printf("dataBytes = %d, lenField = %b, %d, padBytes = %b, %d\n", dataBytes, lenField, lenField, padBytes, padBytes)
	return lenField, padBytes
}
