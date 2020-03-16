package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	const sample = "\xbd\xb2\x3d\xbc\x20\xe2\x8c\x98"

	for _, iRune := range []rune(sample) {
		fmt.Printf("%c ", iRune)
	}
	fmt.Println("\nPrintln:")
	fmt.Println(sample)

	fmt.Println("Byte loop:")
	for i := 0; i < len(sample); i++ {
		fmt.Printf("%x ", sample[i])
	}
	fmt.Printf("\n")

	fmt.Println("Printf with %x:")
	fmt.Printf("%x\n", sample)

	fmt.Println("Printf with % x:")
	fmt.Printf("% x\n", sample)

	fmt.Println("Printf with %q:")
	fmt.Printf("%q\n", sample)

	fmt.Println("Printf with %+q:")
	fmt.Printf("%+q\n", sample)

	fmt.Println("----------------------------------------")
	const placeOfInterest = `⌘`

	fmt.Printf("plain string: ")
	fmt.Printf("%s", placeOfInterest)
	fmt.Printf("\n")

	fmt.Printf("quoted string: ")
	fmt.Printf("%+q", placeOfInterest)
	fmt.Printf("\n")

	fmt.Printf("hex bytes: ")
	for i := 0; i < len(placeOfInterest); i++ {
		fmt.Printf("%x ", placeOfInterest[i])
	}

	for _, ch := range placeOfInterest {
		fmt.Printf("\nUnicode character: %c", ch)
	}
	fmt.Printf("\nThe length of placeOfInterest: %d", len(placeOfInterest))
	fmt.Printf("\n")

	fmt.Println("----------------------------------------")

	const nihongo = "日本語"
	fmt.Println(len(nihongo))
	for index, runeValue := range nihongo {
		fmt.Printf("%#U starts at byte position %d\n", runeValue, index)
	}
	for i, w := 0, 0; i < len(nihongo); i += w {
		runeValue, width := utf8.DecodeRuneInString(nihongo[i:])
		fmt.Printf("%#U starts at byte position %d\n", runeValue, i)
		w = width
	}

	fmt.Println("----------------------------------------")

	a := "Randal"
	for i := 0; i < len(a); i++ {
		fmt.Printf("%x ", a[i])
		fmt.Printf("%c ", a[i])
	}

	fmt.Println("\n----------------------------------------")

	var s string
	s = "中国string"
	r := []rune(s)
	for i := 0; i < len(r); i++ {
		fmt.Printf("%x ", r[i])
	}
	for i := 0; i < len(r); i++ {
		fmt.Printf("%c ", r[i])
	}

	fmt.Println("\n----------------------------------------")

	// for range对字符串进行遍历时，每次获取到的对象都是rune类型的
	s = "知"
	for _, item := range s {
		fmt.Printf("%c", item)
	}
	fmt.Println(len(s))
	for i := 0; i < len(s); i++ {
		fmt.Printf("%c  ", s[i])
	}
	for i := 0; i < len(s); i++ {
		fmt.Printf("%x  ", s[i])
	}

	fmt.Println("\n---------utf8-------------------------------")

	// invalid
	utf8Bytes := []byte{0b11110000, 0b10000000, 0b10000000, 0b10111000}
	// utf8Bytes := []byte{0b11110010, 0b10000010, 0b10000000, 0b10111000}
	fmt.Println(utf8Bytes)
	ru, size := utf8.DecodeRune(utf8Bytes)
	fmt.Println(ru)
	fmt.Printf("%b\n", ru)
	fmt.Println(size)

}

// Println:
// = ⌘
// Byte loop:
// bd b2 3d bc 20 e2 8c 98
// Printf with %x:
// bdb23dbc20e28c98
// Printf with % x:
// bd b2 3d bc 20 e2 8c 98
// Printf with %q:
// "\xbd\xb2=\xbc ⌘"
// Printf with %+q:
// "\xbd\xb2=\xbc \u2318"
// ----------------------------------------
// plain string: ⌘
// quoted string: "\u2318"
// hex bytes: e2 8c 98
