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
	path := "/Users/xsky/go/src/github.com/etcd-io/etcd/infra1.etcd/member/wal/0000000000000000-0000000000000000.wal"
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

	fmt.Println("==============================")
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

	fmt.Println("==============================")
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
