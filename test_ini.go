package main

import (
	"fmt"
	"github.com/go-ini/ini"
	"log"
	"os"
)

func main() {
	x, err := ini.Load("/Users/xsky/.aws/config")
	fmt.Println(x, err)

	y, err := ini.Load(os.Getenv("HOME") + "/.aws/config")
	fmt.Println(y, err)
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
}
