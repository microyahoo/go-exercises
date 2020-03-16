package main

import "fmt"

type person struct {
	name string
	age  int
}

// Note: cannot assign to struct field p["HM"].age in map
type people map[string]person

// type people map[string]*person

func main() {
	p := make(people)
	p["HM"] = person{"Hank McNamara", 39}
	p["HM"].age++
	fmt.Printf("age: %d\n", p["HM"].age)
}
