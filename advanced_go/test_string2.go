package main

import (
	"fmt"
	"reflect"
	"unicode/utf8"
	"unsafe"
)

func main() {
	s := "hello world"
	sHdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	// fmt.Printf("%p, %#v\n", sHdr, sHdr)
	// fmt.Printf("%p, %#v\n", sHdr.Data, sHdr.Data)
	fmt.Println(sHdr.Len)
	fmt.Println(*(*byte)(unsafe.Pointer(sHdr.Data)))
	// fmt.Printf("%p, %#v\n", sHdr, sHdr)
	b := *(*[]byte)(unsafe.Pointer(sHdr))
	// b := *(*[]byte)(unsafe.Pointer(sHdr.Data)) // [signal SIGSEGV: segmentation violation

	// fmt.Printf("%p\n", &b)
	// fmt.Printf("%p, %#v\n", sHdr, sHdr)
	fmt.Println(len(b), cap(b))
	fmt.Println(string(b[:len(s)]))
	fmt.Println(b[:len(s)])
	bs := str2bytes(s)
	fmt.Println("&&&&\n", bs)
	bs2 := string2bytes(s)
	fmt.Println(bs2)
	fmt.Println(bytes2string(bs2))
	// bs2[1] = 108  // signal SIGBUS: bus error code=0x2
	fmt.Println(bs2[1])
	fmt.Println(bytes2string(bs2))

	newbs := []byte{104, 101, 108, 108, 111, 32, 119, 111, 114, 108, 100}
	str := bytes2string(newbs)
	fmt.Println("****\n", str)
	strToB := string2bytes(str)
	strToB[1] = 100
	fmt.Println(bytes2string(newbs))
}

func str2bytes(s string) []byte {
	p := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		p[i] = s[i]
	}
	return p
}

func runes2string(s []int32) string {
	var p []byte
	buf := make([]byte, 3)
	for _, r := range s {
		n := utf8.EncodeRune(buf, r)
		p = append(p, buf[:n]...)
	}
	return string(p)
}

func string2bytes(s string) []byte {
	sHdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
	slHdr := &reflect.SliceHeader{
		Data: sHdr.Data,
		Len:  sHdr.Len,
		Cap:  sHdr.Len,
	}
	return *(*[]byte)(unsafe.Pointer(slHdr))
}

func bytes2string(b []byte) string {
	sHdr := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	slHdr := &reflect.StringHeader{
		Data: sHdr.Data,
		Len:  sHdr.Len,
	}
	return *(*string)(unsafe.Pointer(slHdr))
}
