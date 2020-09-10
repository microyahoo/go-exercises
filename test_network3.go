package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println(interfaceAddrs())
}

func interfaceAddrs() (addrList []string, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, cidr := range addrs {
		// CIDR formats: 172.17.0.1/16 ::1/128
		addr := strings.Split(cidr.String(), "/")[0]
		ip := net.ParseIP(addr)
		if ip != nil {
			if !ip.IsLoopback() {
				addrList = append(addrList, ip.String())
			}
		}
	}

	return addrList, nil
}
