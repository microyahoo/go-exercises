package main

import (
	"fmt"
	// "reflect"
	// "unicode/utf8"
)

// Numbers fundamental to the encoding.
const (
	RuneError = '\uFFFD'     // the "error" Rune or "Unicode replacement character"
	RuneSelf  = 0x80         // characters below Runeself are represented as themselves in a single byte.
	MaxRune   = '\U0010FFFF' // Maximum valid Unicode code point.
	UTFMax    = 4            // maximum number of bytes of a UTF-8 encoded Unicode character.
)

const (
	t1 = 0x00 // 0000 0000
	tx = 0x80 // 1000 0000
	t2 = 0xC0 // 1100 0000
	t3 = 0xE0 // 1110 0000
	t4 = 0xF0 // 1111 0000
	t5 = 0xF8 // 1111 1000

	maskx = 0x3F // 0011 1111
	mask2 = 0x1F // 0001 1111
	mask3 = 0x0F // 0000 1111
	mask4 = 0x07 // 0000 0111

	rune1Max = 1<<7 - 1
	rune2Max = 1<<11 - 1
	rune3Max = 1<<16 - 1

	// The default lowest and highest continuation byte.
	locb = 0x80 // 1000 0000
	hicb = 0xBF // 1011 1111

	// These names of these constants are chosen to give nice alignment in the
	// table below. The first nibble is an index into acceptRanges or F for
	// special one-byte cases. The second nibble is the Rune length or the
	// Status for the special one-byte case.
	xx = 0xF1 // invalid: size 1
	as = 0xF0 // ASCII: size 1
	s1 = 0x02 // accept 0, size 2
	s2 = 0x13 // accept 1, size 3
	s3 = 0x03 // accept 0, size 3
	s4 = 0x23 // accept 2, size 3
	s5 = 0x34 // accept 3, size 4
	s6 = 0x04 // accept 0, size 4
	s7 = 0x44 // accept 4, size 4
)

type acceptRange struct {
	lo uint8 // lowest value for second byte.
	hi uint8 // highest value for second byte.
}

var acceptRanges = [...]acceptRange{
	0: {locb, hicb},
	1: {0xA0, hicb},
	2: {locb, 0x9F},
	3: {0x90, hicb},
	4: {locb, 0x8F},
}

// first is information about the first byte in a UTF-8 sequence.
var first = [256]uint8{
	//   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x00-0x0F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x10-0x1F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x20-0x2F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x30-0x3F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x40-0x4F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x50-0x5F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x60-0x6F
	as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, as, // 0x70-0x7F
	//   1   2   3   4   5   6   7   8   9   A   B   C   D   E   F
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0x80-0x8F
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0x90-0x9F
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0xA0-0xAF
	xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0xB0-0xBF
	xx, xx, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, // 0xC0-0xCF
	s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, s1, // 0xD0-0xDF
	s2, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s3, s4, s3, s3, // 0xE0-0xEF
	s5, s6, s6, s6, s7, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, xx, // 0xF0-0xFF
}

// RuneCountInString is like RuneCount but its input is a string.
func RuneCountInString(s string) (n int) {
	ns := len(s)
	fmt.Println(ns)
	for i := 0; i < ns; n++ {
		c := s[i]
		if c < RuneSelf {
			// ASCII fast path
			i++
			continue
		}
		fmt.Println("c=", c)
		x := first[c]
		fmt.Println("x=", x)
		if x == xx {
			i++ // invalid.
			continue
		}
		size := int(x & 7)
		fmt.Println("size=", size)
		if i+size > ns {
			i++ // Short or invalid.
			continue
		}
		accept := acceptRanges[x>>4]
		fmt.Println("accept: ", accept)
		if c := s[i+1]; c < accept.lo || accept.hi < c {
			size = 1
		} else if size == 2 {
		} else if c := s[i+2]; c < locb || hicb < c {
			size = 1
		} else if size == 3 {
		} else if c := s[i+3]; c < locb || hicb < c {
			size = 1
		}
		i += size
	}
	return n
}

// FullRune ...
func FullRune(p []byte) bool {
	n := len(p)
	if n == 0 {
		return false
	}
	fmt.Println("po=", p[0])
	x := first[p[0]]
	if n >= int(x&7) {
		return true // ASCII, invalid or valid.
	}
	// Must be short or invalid.
	accept := acceptRanges[x>>4]
	if n > 1 && (p[1] < accept.lo || accept.hi < p[1]) {
		return true
	} else if n > 2 && (p[2] < locb || hicb < p[2]) {
		return true
	}
	return false
}

// FullRuneInString is like FullRune but its input is a string.
func FullRuneInString(s string) bool {
	n := len(s)
	if n == 0 {
		return false
	}
	x := first[s[0]]
	fmt.Println("xxx= ", x)
	fmt.Println("x&7= ", x&7)
	if n >= int(x&7) {
		fmt.Println("--------")
		return true // ASCII, invalid, or valid.
	}
	// Must be short or invalid.
	accept := acceptRanges[x>>4]
	if n > 1 && (s[1] < accept.lo || accept.hi < s[1]) {
		fmt.Println("xxxxxx")
		return true
	} else if n > 2 && (s[2] < locb || hicb < s[2]) {
		fmt.Println("eeeee")
		return true
	}
	return false
}

// DecodeRune ...
func DecodeRune(p []byte) (r rune, size int) {
	n := len(p)
	if n < 1 {
		return RuneError, 0
	}
	p0 := p[0]
	fmt.Println("p0 = ", p0)
	x := first[p0]
	fmt.Println("x = ", x)
	if x >= as {
		// The following code simulates an additional check for x == xx and
		// handling the ASCII and invalid cases accordingly. This mask-and-or
		// approach prevents an additional branch.
		mask := rune(x) << 31 >> 31 // Create 0x0000 or 0xFFFF.
		return rune(p[0])&^mask | RuneError&mask, 1
	}
	sz := int(x & 7)
	fmt.Println("size = ", sz)
	accept := acceptRanges[x>>4]
	fmt.Println("accept = ", accept)
	if n < sz {
		return RuneError, 1
	}
	b1 := p[1]
	fmt.Printf("b1 = %d\n", b1)
	if b1 < accept.lo || accept.hi < b1 {
		return RuneError, 1
	}
	if sz <= 2 { // <= instead of == to help the compiler eliminate some bounds checks
		return rune(p0&mask2)<<6 | rune(b1&maskx), 2
	}
	b2 := p[2]
	fmt.Printf("b2 = %d\n", b2)
	if b2 < locb || hicb < b2 {
		return RuneError, 1
	}
	if sz <= 3 {
		return rune(p0&mask3)<<12 | rune(b1&maskx)<<6 | rune(b2&maskx), 3
	}
	b3 := p[3]
	if b3 < locb || hicb < b3 {
		return RuneError, 1
	}
	return rune(p0&mask4)<<18 | rune(b1&maskx)<<12 | rune(b2&maskx)<<6 | rune(b3&maskx), 4
}

func main() {
	fmt.Printf("acceptRanges = %T\n", acceptRanges)
	utf8Bytes := []byte{
		0b11110000,
		0b10000000,
		0b10000000,
		0b10111000,
	}
	// utf8Bytes := []byte{0b11110010, 0b10000010, 0b10000000, 0b10111000}
	fmt.Println(utf8Bytes)
	ru, size := DecodeRune(utf8Bytes)
	fmt.Println(ru)
	fmt.Printf("%b\n", ru)
	fmt.Println(size)

	fmt.Println("\n----------------------------------------")

	//fmt.Println(reflect.TypeOf(acceptRanges))
	//str := "Hello, 钢铁侠"
	//fmt.Println(FullRuneInString(`\ubbbbbbb`))
	//fmt.Println(FullRune([]byte(str)))
	//fmt.Println(utf8.RuneCount([]byte(str)))
	//fmt.Println(str)
	//for i := 0; i < len(str); i++ {
	//	fmt.Println(str[i])
	//}
	//fmt.Println([]byte(str))
	//for _, s := range str {
	//	fmt.Println(s)
	//}
	//fmt.Println(reflect.TypeOf([]rune(str)[4]))
	//fmt.Println([]rune(str))
	//fmt.Println([]int32(str))
	//fmt.Println(utf8.RuneCountInString(str))
	////fmt.Println(first[uint8(str[6])])
	////accept := acceptRanges[4]
	//fmt.Println(RuneCountInString(str))
	//fmt.Println(utf8.ValidString(str))
}
