package main

import (
	"fmt"
	"strconv"
	"strings"
)

func main() {
	var input []int

	var inputVar string
	fmt.Scan(&inputVar)

	inputSlice := strings.Split(inputVar, ",")
	for _, in := range inputSlice {
		i, err := strconv.Atoi(in)
		if err != nil {
			return
		}
		input = append(input, i)
	}
	// fmt.Println(input)

	low := 0
	high := len(input) - 1
	gap := 1

	if len(input) == 1 {
		fmt.Println(input[0])
		return
	}
	i := low
	j := low + 1
	var output []string
	for j <= high {
		if input[j]-input[i] == gap {
			fmt.Println(j, "---")
			if j == high {
				first := strconv.Itoa(input[i])
				last := strconv.Itoa(input[j])
				if j-i >= 1 {
					output = append(output, fmt.Sprintf("%s-%s", first, last))
				} else {
					output = append(output, first)
				}
			}
			gap++
			j++
		} else {
			first := strconv.Itoa(input[i])
			last := strconv.Itoa(input[j-1])
			if j-i > 1 {
				output = append(output, fmt.Sprintf("%s-%s", first, last))
			} else {
				output = append(output, first)
			}
			if j == high {
				first := strconv.Itoa(input[i])
				last := strconv.Itoa(input[j])
				if j-i > 1 {
					output = append(output, fmt.Sprintf("%s-%s", first, last))
				} else {
					output = append(output, first)
				}
			}
			gap = 1
			i = j
			j++
			fmt.Println(i, j, "===")
		}
	}

	fmt.Println(strings.Join(output, ","))
}
