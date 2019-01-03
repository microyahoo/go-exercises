package main

import (
	"fmt"
	"net"
)

func main() {
	ip := net.ParseIP("10.252.1.129")
	// mask := 32
	vg_network := "10.252.1.3/24"
	ipx, ipNet, _ := net.ParseCIDR(vg_network)
	fmt.Println(ipx)
	fmt.Println(ipNet)

	if ipNet.Contains(ip) {
		fmt.Printf("%v contains %v\n", ipNet, ip)
	} else {
		fmt.Printf("%v not contain %v\n", ipNet, ip)
	}

}
