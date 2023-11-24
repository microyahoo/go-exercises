package main

import (
	"fmt"
	"strings"
)

func main() {
	fmt.Println(strings.Split("vim-go", "/"))
	fmt.Println(strings.Split("", "/"))
	splits := strings.Split("/dev/mapper/ceph--c1491c4a--27a8--4f05--b679--931cb7555a95-osd--block--265d9594--0334--421e--af4c--d960afff5beb", "/")
	fmt.Println(splits, len(splits))
	last := splits[len(splits)-1]
	fmt.Println(last)
	// lastSplits := strings.Split(last, "--")
	// fmt.Println(lastSplits, len(lastSplits))
	newLast := strings.ReplaceAll(last, "--", "#")
	fmt.Println(newLast)
	newLastSplits := strings.Split(newLast, "-")
	fmt.Println(newLastSplits)
	var vg, lv string
	if len(newLastSplits) == 2 {
		vg = strings.ReplaceAll(newLastSplits[0], "#", "-")
		lv = strings.ReplaceAll(newLastSplits[1], "#", "-")
	}
	fmt.Println(vg)
	fmt.Println(lv)
}
