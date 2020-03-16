package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	minIP, maxIP := getCidrIPRange("100.111.111.111/22")
	fmt.Println("CIDR最小IP：", minIP, " CIDR最大IP：", maxIP)
	fmt.Println("掩码：", getCidrIPMask(22))
	fmt.Println("主机数量", getCidrHostNum(22))
}

func getCidrIPRange(cidr string) (string, string) {
	ip := strings.Split(cidr, "/")[0]
	ipSegs := strings.Split(ip, ".")
	maskLen, _ := strconv.Atoi(strings.Split(cidr, "/")[1])
	seg3MinIP, seg3MaxIP := getIPSeg3Range(ipSegs, maskLen)
	seg4MinIP, seg4MaxIP := getIPSeg4Range(ipSegs, maskLen)
	ipPrefix := ipSegs[0] + "." + ipSegs[1] + "."

	return ipPrefix + strconv.Itoa(seg3MinIP) + "." + strconv.Itoa(seg4MinIP),
		ipPrefix + strconv.Itoa(seg3MaxIP) + "." + strconv.Itoa(seg4MaxIP)
}

//计算得到CIDR地址范围内可拥有的主机数量
func getCidrHostNum(maskLen int) uint {
	cidrIPNum := uint(0)
	var i = uint(32 - maskLen - 1)
	// var i uint = uint(32 - maskLen - 1)
	for ; i >= 1; i-- {
		cidrIPNum += 1 << i
	}
	return cidrIPNum
}

//获取Cidr的掩码
func getCidrIPMask(maskLen int) string {
	// ^uint32(0)二进制为32个比特1，通过向左位移，得到CIDR掩码的二进制
	cidrMask := ^uint32(0) << uint(32-maskLen)
	fmt.Println(fmt.Sprintf("%b \n", cidrMask))
	//计算CIDR掩码的四个片段，将想要得到的片段移动到内存最低8位后，将其强转为8位整型，从而得到
	cidrMaskSeg1 := uint8(cidrMask >> 24)
	cidrMaskSeg2 := uint8(cidrMask >> 16)
	cidrMaskSeg3 := uint8(cidrMask >> 8)
	cidrMaskSeg4 := uint8(cidrMask & uint32(255))

	return fmt.Sprint(cidrMaskSeg1) + "." + fmt.Sprint(cidrMaskSeg2) + "." + fmt.Sprint(cidrMaskSeg3) + "." + fmt.Sprint(cidrMaskSeg4)
}

//得到第三段IP的区间（第一片段.第二片段.第三片段.第四片段）
func getIPSeg3Range(ipSegs []string, maskLen int) (int, int) {
	if maskLen > 24 {
		segIP, _ := strconv.Atoi(ipSegs[2])
		return segIP, segIP
	}
	ipSeg, _ := strconv.Atoi(ipSegs[2])
	return getIPSegRange(uint8(ipSeg), uint8(24-maskLen))
}

//得到第四段IP的区间（第一片段.第二片段.第三片段.第四片段）
func getIPSeg4Range(ipSegs []string, maskLen int) (int, int) {
	ipSeg, _ := strconv.Atoi(ipSegs[3])
	segMinIP, segMaxIP := getIPSegRange(uint8(ipSeg), uint8(32-maskLen))
	return segMinIP + 1, segMaxIP
}

//根据用户输入的基础IP地址和CIDR掩码计算一个IP片段的区间
func getIPSegRange(userSegIP, offset uint8) (int, int) {
	var ipSegMax uint8 = 255
	netSegIP := ipSegMax << offset
	segMinIP := netSegIP & userSegIP
	segMaxIP := userSegIP&(255<<offset) | ^(255 << offset)
	return int(segMinIP), int(segMaxIP)
}
