package main

import (
	"fmt"

	"github.com/containerd/cgroups/v3"
)

func isCgroupV2() bool {
	return cgroups.Mode() == cgroups.Unified
}

func main() {
	if isCgroupV2() {
		fmt.Println("System is using cgroup v2 (Unified)")
	} else {
		fmt.Println("System is using cgroup v1 (Legacy or Hybrid)")
	}
}
