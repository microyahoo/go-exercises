package main

import (
	"fmt"
	"sort"
)

func main() {
	testDifferent1()
	fmt.Println("========================")
	testDifferent2()
}

func testDifferent1() {
	var allOsdIDs, excludeOsdIDs []uint64
	var i uint64
	for i = 0; i < 303; i++ {
		allOsdIDs = append(allOsdIDs, i)
	}
	for i = 100; i < 123; i++ {
		excludeOsdIDs = append(excludeOsdIDs, i)
	}
	different := func(list1, list2 []uint64) (result []uint64) {
		for _, i := range list1 {
			existing := false
			for _, j := range list2 {
				if i == j {
					existing = true
					break
				}
			}
			if !existing {
				result = append(result, i)
			}
		}
		return result
	}
	remainingOsdIDs := different(allOsdIDs, excludeOsdIDs)
	fmt.Println(remainingOsdIDs)
	osds := make([]uint64, 0)
	listOsdsByOsdIDs := func(osdIDs []uint64) (osds []uint64) {
		for _, id := range osdIDs {
			osds = append(osds, id+1000)
		}
		return osds
	}
	for len(remainingOsdIDs) > 100 {
		fmt.Println("++\n", len(remainingOsdIDs))
		osdSlice := listOsdsByOsdIDs(remainingOsdIDs[:100])
		osds = append(osds, osdSlice...)
		remainingOsdIDs = remainingOsdIDs[100:]
	}
	if len(remainingOsdIDs) > 0 {
		fmt.Println("**\n", len(remainingOsdIDs))
		osdSlice := listOsdsByOsdIDs(remainingOsdIDs)
		osds = append(osds, osdSlice...)
	}
	fmt.Println("\n", osds)
	fmt.Println("\n", len(osds))
}

func testDifferent2() {
	var allOsdIDs, excludeOsdIDs []uint64
	var i uint64
	for i = 0; i < 303; i++ {
		allOsdIDs = append(allOsdIDs, i)
	}
	allOsdIDs = append(allOsdIDs, 1000000)
	for i = 100; i < 123; i++ {
		excludeOsdIDs = append(excludeOsdIDs, i)
	}
	excludeOsdIDs = append(excludeOsdIDs, 12390000)
	different := func(list1, list2 []uint64) (result []uint64) {
		map2 := make(map[uint64]struct{}, len(list2))
		for _, j := range list2 {
			if _, ok := map2[j]; !ok {
				map2[j] = struct{}{}
			}
		}
		for _, i := range list1 {
			if _, ok := map2[i]; !ok {
				result = append(result, i)
			}
		}
		sort.Slice(result, func(i, j int) bool { return result[i] < result[j] })
		return result
	}
	remainingOsdIDs := different(allOsdIDs, excludeOsdIDs)
	fmt.Println(remainingOsdIDs)
	osds := make([]uint64, 0)
	listOsdsByOsdIDs := func(osdIDs []uint64) (osds []uint64) {
		for _, id := range osdIDs {
			osds = append(osds, id+1000)
		}
		return osds
	}
	min := func(i, j int) bool {
		return i < j
	}
	step := 100
	for len(remainingOsdIDs) > 0 {
		fmt.Println("++\n", len(remainingOsdIDs))
		if min(len(remainingOsdIDs), step) {
			step = len(remainingOsdIDs)
		}
		osdSlice := listOsdsByOsdIDs(remainingOsdIDs[:step])
		osds = append(osds, osdSlice...)
		remainingOsdIDs = remainingOsdIDs[step:]
	}
	// if len(remainingOsdIDs) > 0 {
	// 	fmt.Println("**\n", len(remainingOsdIDs))
	// 	osdSlice := listOsdsByOsdIDs(remainingOsdIDs)
	// 	osds = append(osds, osdSlice...)
	// }
	fmt.Println("\n", osds)
	fmt.Println("\n", len(osds))

	var tmps []uint64
	func() {
		count := len(osds)
		tmps = osds[:count]
		osds = osds[count:]
	}()
	fmt.Println("===========*============")
	fmt.Println("\n", osds)
	fmt.Println("\n", len(osds))
	fmt.Println("\n", tmps)
	fmt.Println("\n", len(tmps))
}
