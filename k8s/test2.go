package main

import (
	"fmt"
	"os"
	"syscall"
)

func main() {
	fmt.Println(getFileInode("/root/jfs"))
	fmt.Println(getFileInode("/etc"))
	fmt.Println(getFileInode("/"))
}

func getFileInode(path string) (uint64, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	if sst, ok := fi.Sys().(*syscall.Stat_t); ok {
		fmt.Printf("stat %+v\n", sst)
		return sst.Ino, nil
	}
	return 0, nil
}
