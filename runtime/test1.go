package main

import (
	"fmt"
	"io/fs"
	"os"
	"sort"
	"time"

	"github.com/microyahoo/go-exercises/runtime/outer"
	gopsutildisk "github.com/shirou/gopsutil/v3/disk"
	"golang.org/x/sys/unix"
)

// https://www.pixelstech.net/article/1649596852-The-magic-of-go%3Alinkname
func main() {
	time.Sleep(time.Second)
	outer.World()

	fmt.Println(1 / 2)
	fmt.Println(1 / (1.0 * 2))

	logDir := "/var/lib/rook/rook-ceph/log"
	entries, err := os.ReadDir(logDir)
	if err != nil {
		panic(err)
	}
	var infos []fs.FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fileInfo, err := entry.Info()
		if err != nil {
			panic(err)
		}
		infos = append(infos, fileInfo)
	}
	// sort by modification time
	sort.Slice(infos, func(i, j int) bool {
		return infos[i].ModTime().Before(infos[j].ModTime())
	})
	for _, info := range infos {
		fmt.Printf("%s(%d)\n", info.Name(), info.Size())
	}

	_, capacity, usage, _, _, _, _ := fsInfo("/")
	fmt.Println(usage, capacity)
	fmt.Println(usage / (1.0 * capacity))
	fmt.Println(float64(usage / capacity))
	fmt.Println(float64(usage / (1.0 * capacity)))
	fmt.Println(float64(1.0 * usage / capacity))
	fmt.Println(float64(usage) / float64(capacity))
	fmt.Println(1.0 * usage)
	fmt.Printf("%T\n", 1.0*usage)
	fmt.Printf("%T\n", usage*1.0)

	usageStat, err := gopsutildisk.Usage("/")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", usageStat)
}

// FsInfo linux returns (available bytes, byte capacity, byte usage, total inodes, inodes free, inode usage, error)
// for the filesystem that path resides upon.
func fsInfo(path string) (int64, int64, int64, int64, int64, int64, error) {
	statfs := &unix.Statfs_t{}
	err := unix.Statfs(path, statfs)
	if err != nil {
		return 0, 0, 0, 0, 0, 0, err
	}

	// Available is blocks available * fragment size
	available := int64(statfs.Bavail) * int64(statfs.Bsize)

	// Capacity is total block count * fragment size
	capacity := int64(statfs.Blocks) * int64(statfs.Bsize)

	// Usage is block being used * fragment size (aka block size).
	usage := (int64(statfs.Blocks) - int64(statfs.Bfree)) * int64(statfs.Bsize)

	inodes := int64(statfs.Files)
	inodesFree := int64(statfs.Ffree)
	inodesUsed := inodes - inodesFree

	return available, capacity, usage, inodes, inodesFree, inodesUsed, nil
}
