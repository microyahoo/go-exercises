package main

import (
	"fmt"
	"regexp"
)

var (
	temperatureRegex = regexp.MustCompile(`(\s+)?(?P<temp>\d+)(\s+\(.*\))?`)
)

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
}
