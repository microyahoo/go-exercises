package main

import "fmt"

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
