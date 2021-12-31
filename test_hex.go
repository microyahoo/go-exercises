package main

import (
	"fmt"
	"os"
)

func main() {
	indices := []int64{0, 0xc41e9b, 0xc445ac, 0xc46cbd, 0xc493ce, 0xc4badf}
	if len(indices) < 2 {
		fmt.Println("The length of indices is less than 2")
	} else {
		i := 1
		for i < len(indices) {
			index0 := indices[i-1]
			index1 := indices[i]
			fmt.Printf("The diff between indices[%x, %x] is %d\n", index0, index1, index1-index0)
			i++
		}
	}

	fmt.Printf("14488752894282421173 = %x\n", uint64(14488752894282421173))
	fmt.Printf("5032027080688003127 = %x\n", uint64(5032027080688003127))
	fmt.Printf("18249187646912138824 = %x\n", uint64(18249187646912138824))
	fmt.Printf("10501334649042878790 = %x\n", uint64(10501334649042878790))
	fmt.Printf("9372538179322589801 = %x\n", uint64(9372538179322589801))
	fmt.Printf("7044390442471271841 = %x\n", uint64(7044390442471271841))
	fmt.Printf("0xf4e = %d\n", uint64(0xf4e))
	fmt.Printf("13939 = %x\n", uint64(13939))
	fmt.Printf("%d = 0x%x\n", uint64(0x3868), 0x3868)
	fmt.Printf("%d = 0x%x\n", uint64(0x3a5d), 0x3a5d)
	fmt.Printf("%d = 0x%x\n", uint64(0x59ad), 0x59ad)
	fmt.Printf("%d = 0x%x\n", uint64(0x5ba2), 0x5ba2)
	fmt.Printf("%d = 0x%x\n", uint64(0x5d82), 0x5d82)
	fmt.Printf("%d = 0x%x\n", uint64(0x32697dff1a127f04), 0x32697dff1a127f04)
	fmt.Printf("%d = 0x%x\n", uint64(0x32697dfb4680e707), 0x32697dfb4680e707)

	os.Rename("1", "2")
}

// (.venv) /root ☞ tree -h /var/lib/etcd/
// /var/lib/etcd/
// └── [  29]  member
//     ├── [ 246]  snap
//     │         ├── [ 12K]  000000000000000e-0000000000c41e9b.snap
//     │         ├── [ 12K]  000000000000000e-0000000000c445ac.snap
//     │         ├── [ 12K]  000000000000000e-0000000000c46cbd.snap
//     │         ├── [ 12K]  000000000000000e-0000000000c493ce.snap
//     │         ├── [ 12K]  000000000000000e-0000000000c4badf.snap
//     │         └── [497M]  db
//     └── [ 244]  wal
//         ├── [ 61M]  0000000000000087-0000000000be47f8.wal
//         ├── [ 61M]  0000000000000088-0000000000bfb854.wal
//         ├── [ 61M]  0000000000000089-0000000000c128b5.wal
//         ├── [ 61M]  000000000000008a-0000000000c29920.wal
//         ├── [ 61M]  000000000000008b-0000000000c4094d.wal
//         └── [ 61M]  1.tmp
//
// 3 directories, 12 files
