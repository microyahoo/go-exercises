package main

import "fmt"

func hello1(i *int) int {
	defer func() {
		*i = 19
	}()
	return *i
}

func hello2(i *int) (j int) {
	defer func() {
		j = 19
	}()
	j = *i
	return j
}

func hello3(i *int) (j int) {
	defer func() {
		j = 19
	}()
	return *i
}

func main() {
	i := 10
	j := hello1(&i)
	fmt.Println(i, j)

	i = 10
	j = hello2(&i)
	fmt.Println(i, j)

	i = 10
	j = hello3(&i)
	fmt.Println(i, j)
}
