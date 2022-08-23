package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

func main() {
	buf := new(bytes.Buffer)
	var num byte = 250
	err := binary.Write(buf, binary.LittleEndian, num)
	if err != nil {
		panic(fmt.Sprintf("binary.Write failed: %s", err))
	}
	bytes := buf.Bytes()
	fmt.Printf("%b\n", bytes)

	fmt.Printf("250 = %b\n", uint32(bytes[0]))
	fmt.Printf("uint32(250) = %b\n", uint32(bytes[0]))
	fmt.Printf("uint32(int8(250)) = %b\n", uint32(int8(bytes[0])))
}
