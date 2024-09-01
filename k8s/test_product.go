package main

import (
	"fmt"
)

var output [][]string

func product(arr [][]string, index int, cur_res []string) {
	if index == len(arr) {
		res := make([]string, 0, len(cur_res))
		for _, ele := range cur_res {
			res = append(res, ele)
		}
		output = append(output, res)
		return
	}
	for j := 0; j < len(arr[index]); j++ {
		cur_res = append(cur_res, arr[index][j])
		product(arr, index+1, cur_res)
		cur_res = cur_res[:len(cur_res)-1]
	}
}

func main() {
	output = make([][]string, 0)
	input := make([][]string, 4, 4)
	input[0] = make([]string, 2, 2)
	input[1] = make([]string, 2, 2)
	input[2] = make([]string, 3, 3)
	input[3] = make([]string, 1, 1)
	input[0][0] = "a1"
	input[0][1] = "a2"
	input[1][0] = "b1"
	input[1][1] = "b2"
	input[2][0] = "c1"
	input[2][1] = "c2"
	input[2][2] = "c3"
	input[3][0] = "d1"
	product(input, 0, make([]string, 0))
	fmt.Println(output)
}
