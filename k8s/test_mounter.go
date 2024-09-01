package main

import (
	"fmt"

	mount "k8s.io/mount-utils"
)

func main() {
	var (
		mounter mount.Interface = mount.NewWithoutSystemd("")
		p                       = "/var/lib/kubelet/pods/c9d3dfb0-3b18-40a2-9d99-ed09261d170d/volumes/kubernetes.io~csi/pvc-dce9694c-2e96-4324-990b-e6dcb21c87f7/mount"
	)
	notMnt, err := mounter.IsLikelyNotMountPoint(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(notMnt)
	p = "/mnt/xxx"
	notMnt, err = mounter.IsLikelyNotMountPoint(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(notMnt)
	p = "/tmp/b" // mkdir /tmp/a /tmp/b; mount --bind /tmp/a /tmp/b
	notMnt, err = mounter.IsLikelyNotMountPoint(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(notMnt)
}
