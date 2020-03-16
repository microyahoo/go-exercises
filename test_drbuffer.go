package main

import (
	"fmt"

	"github.com/esdb/drbuffer"
)

func main() {
	buffer, err := drbuffer.Open("/tmp/drbuffer", 1 /*in kb*/)
	if err != nil {
		return
	}
	defer buffer.Close()
	// push one packet ([]byte)
	// must ensure the packet pushed do not exceed 65535 byte
	buffer.PushOne([]byte("Hello"))
	// batch push multiple packets
	buffer.PushN([][]byte{
		[]byte("A"),
		[]byte("B"),
	})
	// if nothing to pop, packet will be nil
	packet := buffer.PopOne()
	fmt.Println(string(packet))
	// packets is array of []byte
	packets := buffer.PopN(2)
	fmt.Println(packets)
}
