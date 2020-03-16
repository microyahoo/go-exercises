package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func main() {
	out, err := ioutil.ReadFile("/tmp/a")
	if err != nil {
		fmt.Println(err)
	}
	str := strings.TrimSpace(string(out))
	if str == "" {
		fmt.Println("file is null")
	} else {
		fmt.Println(str)
	}
}
