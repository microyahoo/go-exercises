package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println("vim-go")
	groups := []string{
		"osd",
		"pool",
		"disk",
		"lun",
	}
	fmt.Println(strings.Join(groups, "',"))
	fmt.Println(addHexPrefix("hello"))
}

func addHexPrefix(s string) string {
	ns := make([]byte, len(s)*2)
	for i := 0; i < len(s); i += 2 {
		ns[i*2] = '\\'
		ns[i*2+1] = 'x'
		ns[i*2+2] = s[i]
		ns[i*2+3] = s[i+1]
	}
	return string(ns)
}
