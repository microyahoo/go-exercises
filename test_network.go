package main

import (
	"fmt"
	"net"
	"strings"
)

// HandleIQNCode handle iqn code
func HandleIQNCode(code string) string {
	if strings.IndexByte(code, ':') != -1 && net.ParseIP(code) != nil {
		ip := net.ParseIP(code).String()
		if strings.IndexByte(ip, ':') != -1 {
			return "[" + ip + "]"
		}
		return ip
	}
	return code
}

// RecoverIQNCode recovers iqn code
func RecoverIQNCode(code string) string {
	start := strings.IndexByte(code, '[')
	end := strings.IndexByte(code, ']')
	if start == 0 && end == len(code)-1 && net.ParseIP(code[1:end]) != nil {
		return code[1:end]
	}
	return code
}

func main() {
	// ip := net.ParseIP("10.252.1.128")
	ip := net.ParseIP("172.16.149.127")
	// ip := net.ParseIP("172.16.149.128")
	// ip := net.ParseIP("172.16.149.129")
	// ip := net.ParseIP("172.16.149.120")
	handleIP := HandleIQNCode(ip.String())
	fmt.Printf("%s => %s => %s\n", ip.String(), handleIP, RecoverIQNCode(handleIP))
	// mask := 32
	// vgNetwork := "172.16.149.65/26"
	// vgNetwork := "172.16.149.0/26"
	vgNetwork := "192.168.129.53/14"
	handleIP = HandleIQNCode(vgNetwork)
	fmt.Printf("%s => %s => %s\n", vgNetwork, handleIP, RecoverIQNCode(handleIP))
	ipx, ipNet, _ := net.ParseCIDR(vgNetwork)
	fmt.Println(ipx)
	fmt.Println(ipNet)

	if ipNet.Contains(ip) {
		fmt.Printf("%v contains %v\n", ipNet, ip)
	} else {
		fmt.Printf("%v not contain %v\n", ipNet, ip)
	}

	ipv6 := "::FFFF:192.168.0.1"
	handleIPv6 := HandleIQNCode(ipv6)
	fmt.Printf("%s => %s => %s\n", ipv6, handleIPv6, RecoverIQNCode(handleIPv6))
	fmt.Println(net.ParseIP(ipv6))
	ipv6 = "::192.168.0.1"
	handleIPv6 = HandleIQNCode(ipv6)
	fmt.Printf("%s => %s => %s\n", ipv6, handleIPv6, RecoverIQNCode(handleIPv6))
	fmt.Println(net.ParseIP(ipv6))
	ipv6 = "::1"
	handleIPv6 = HandleIQNCode(ipv6)
	fmt.Printf("%s => %s => %s\n", ipv6, handleIPv6, RecoverIQNCode(handleIPv6))
	fmt.Println(net.ParseIP(ipv6))
	ipv6 = "2001:0DB8:0000:0023:0008:0800:200C:417A"
	handleIPv6 = HandleIQNCode(ipv6)
	fmt.Printf("%s => %s => %s\n", ipv6, handleIPv6, RecoverIQNCode(handleIPv6))
	fmt.Println(net.ParseIP(ipv6))
	ipv6 = "[2001:0DB8:0000:0023:0008:0800:200C:417B]"
	handleIPv6 = HandleIQNCode(ipv6)
	fmt.Printf("%s => %s => %s\n", ipv6, handleIPv6, RecoverIQNCode(handleIPv6))
	fmt.Println(net.ParseIP(ipv6))
	ipv6 = "[]"
	handleIPv6 = HandleIQNCode(ipv6)
	fmt.Printf("%s => %s => %s\n", ipv6, handleIPv6, RecoverIQNCode(handleIPv6))
	fmt.Println(net.ParseIP(ipv6))

	wwn := "21:00:00:0e:1e:11:60:50"
	handleWWN := HandleIQNCode(wwn)
	fmt.Printf("%s => %s => %s\n", wwn, handleWWN, RecoverIQNCode(handleWWN))
	fmt.Printf("%v\n", net.ParseIP(wwn))
	fmt.Println(net.ParseIP(wwn))

	addrs, err := net.InterfaceAddrs()
	fmt.Println(addrs, err)
}
