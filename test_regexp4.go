package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var (
	temperatureRegex = regexp.MustCompile(`(\s+)?(?P<temp>\d+)(\s+\(.*\))?`)
	vridRegexp       = regexp.MustCompile("^.*, vrid (\\d+),.*$")
)

// ContainsInt check if array contains an int value
func ContainsInt(slice []int, value int) bool {
	for _, i := range slice {
		if i == value {
			return true
		}
	}
	return false
}

func parseVRIDs(bs []byte) (vrids []int) {
	scanner := bufio.NewScanner(bytes.NewReader(bs))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		strs := vridRegexp.FindStringSubmatch(line)
		if len(strs) != 2 {
			continue
		}
		fmt.Println(strs[0])
		fmt.Println(strs[1])
		fmt.Println("==============")
		vrid, e := strconv.Atoi(strs[1])
		if e != nil {
			log.Printf("invalid vrrid: %s\n", strs[1])
			continue
		}
		if !ContainsInt(vrids, vrid) {
			vrids = append(vrids, vrid)
		}
	}
	slice := sort.IntSlice(vrids)
	slice.Sort()
	return []int(slice)
}

func main() {
	temp1 := "40"
	temp2 := "42 (0 18 0 0 0)"
	temp3 := "  40"
	temp4 := "  40     "
	temp5 := "40      "

	var (
		matches []string
	)
	if temperatureRegex.MatchString(temp1) {
		matches = temperatureRegex.FindStringSubmatch(temp1)
	}
	fmt.Println(matches)

	if temperatureRegex.MatchString(temp2) {
		matches = temperatureRegex.FindStringSubmatch(temp2)
	}
	fmt.Println(matches)

	if temperatureRegex.MatchString(temp3) {
		matches = temperatureRegex.FindStringSubmatch(temp3)
	}
	fmt.Println(matches)

	if temperatureRegex.MatchString(temp4) {
		matches = temperatureRegex.FindStringSubmatch(temp4)
	}
	fmt.Println(matches)

	if temperatureRegex.MatchString(temp5) {
		matches = temperatureRegex.FindStringSubmatch(temp5)
	}
	fmt.Println(matches)
	for i, name := range temperatureRegex.SubexpNames() {
		value := matches[i]
		fmt.Println("i=", i)
		fmt.Println("name=", name)
		fmt.Println("value=", value)
	}

	vridOutput := `16:40:38.598919 IP 172.25.5.75 > 224.0.0.18: VRRPv2, Advertisement, vrid 3, prio 101, authtype simple, intvl 1s, length 20
16:40:38.599226 IP 172.25.5.70 > 224.0.0.18: VRRPv2, Advertisement, vrid 2, prio 101, authtype simple, intvl 1s, length 20
16:40:38.600483 IP 172.25.5.110 > 224.0.0.18: VRRPv2, Advertisement, vrid 255, prio 101, authtype simple, intvl 1s, length 20
16:40:38.904540 IP 172.25.7.152 > 224.0.0.18: VRRPv2, Advertisement, vrid 254, prio 100, authtype simple, intvl 1s, length 20
16:40:39.045179 IP 172.25.5.253 > 224.0.0.18: VRRPv3, Advertisement, vrid 1, prio 120, intvl 100cs, length 12
16:40:39.599326 IP 172.25.5.75 > 224.0.0.18: VRRPv2, Advertisement, vrid 3, prio 101, authtype simple, intvl 1s, length 20
16:40:39.599349 IP 172.25.5.70 > 224.0.0.18: VRRPv2, Advertisement, vrid 2, prio 101, authtype simple, intvl 1s, length 20  
16:40:39.600559 IP 172.25.5.110 > 224.0.0.18: VRRPv2, Advertisement, vrid 255, prio 101, authtype simple, intvl 1s, length 20
16:40:39.905592 IP 172.25.7.152 > 224.0.0.18: VRRPv2, Advertisement, vrid 254, prio 100, authtype simple, intvl 1s, length 20
16:40:40.045174 IP 172.25.5.253 > 224.0.0.18: VRRPv3, Advertisement, vrid 1, prio 120, intvl 100cs, length 12
16:40:40.599705 IP 172.25.5.70 > 224.0.0.18: VRRPv2, Advertisement, vrid 2, prio 101, authtype simple, intvl 1s, length 20
16:40:40.600654 IP 172.25.5.110 > 224.0.0.18: VRRPv2, Advertisement, vrid 255, prio 101, authtype simple, intvl 1s, length 20
16:40:40.601749 IP 172.25.5.75 > 224.0.0.18: VRRPv2, Advertisement, vrid 3, prio 101, authtype simple, intvl 1s, length 20
16:40:40.906027 IP 172.25.7.152 > 224.0.0.18: VRRPv2, Advertisement, vrid 254, prio 100, authtype simple, intvl 1s, length 20
16:40:41.045186 IP 172.25.5.253 > 224.0.0.18: VRRPv3, Advertisement, vrid 1, prio 120, intvl 100cs, length 12
16:40:41.600311 IP 172.25.5.70 > 224.0.0.18: VRRPv2, Advertisement, vrid 2, prio 101, authtype simple, intvl 1s, length 20
16:40:41.601040 IP 172.25.5.110 > 224.0.0.18: VRRPv2, Advertisement, vrid 255, prio 101, authtype simple, intvl 1s, length 20
16:40:41.601838 IP 172.25.5.75 > 224.0.0.18: VRRPv2, Advertisement, vrid 3, prio 101, authtype simple, intvl 1s, length 20
16:40:41.906396 IP 172.25.7.152 > 224.0.0.18: VRRPv2, Advertisement, vrid 254, prio 100, authtype simple, intvl 1s, length 20
16:40:42.045208 IP 172.25.5.253 > 224.0.0.18: VRRPv3, Advertisement, vrid 1, prio 120, intvl 100cs, length 12
16:40:42.600916 IP 172.25.5.70 > 224.0.0.18: VRRPv2, Advertisement, vrid 2, prio 101, authtype simple, intvl 1s, length 20  
16:40:42.601919 IP 172.25.5.75 > 224.0.0.18: VRRPv2, Advertisement, vrid 3, prio 101, authtype simple, intvl 1s, length 20
16:40:42.603026 IP 172.25.5.110 > 224.0.0.18: VRRPv2, Advertisement, vrid 255, prio 101, authtype simple, intvl 1s, length 20
16:40:42.906776 IP 172.25.7.152 > 224.0.0.18: VRRPv2, Advertisement, vrid 254, prio 100, authtype simple, intvl 1s, length 20
16:40:43.045278 IP 172.25.5.253 > 224.0.0.18: VRRPv3, Advertisement, vrid 1, prio 120, intvl 100cs, length 12
16:40:43.601534 IP 172.25.5.70 > 224.0.0.18: VRRPv2, Advertisement, vrid 2, prio 101, authtype simple, intvl 1s, length 20
16:40:43.601998 IP 172.25.5.75 > 224.0.0.18: VRRPv2, Advertisement, vrid 3, prio 101, authtype simple, intvl 1s, length 20
16:40:43.603103 IP 172.25.5.110 > 224.0.0.18: VRRPv2, Advertisement, vrid 255, prio 101, authtype simple, intvl 1s, length 20
16:40:43.907137 IP 172.25.7.152 > 224.0.0.18: VRRPv2, Advertisement, vrid 254, prio 100, authtype simple, intvl 1s, length 20
`
	fmt.Println(parseVRIDs([]byte(vridOutput)))
}
