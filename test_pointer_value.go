package main

import "fmt"

type car struct {
	name  string
	color string
}

func (c *car) SetName01(s string) {
	fmt.Printf("SetName01: car address: %p\n", c)
	c.name = s
}

func (c car) SetName02(s string) {
	fmt.Printf("SetName02: car address: %p\n", &c)
	c.name = s
}

func testxxx() (cars []*car) {
	return cars
}

func main() {
	greatWall := &car{
		name:  "Great Wall",
		color: "white",
	}

	fmt.Printf("car address: %p\n", greatWall)

	fmt.Println(greatWall.name)
	greatWall.SetName01("foo")
	fmt.Println(greatWall.name)

	greatWall.SetName02("bar")
	fmt.Println(greatWall.name)

	greatWall.SetName02("test")
	fmt.Println(greatWall.name)
	fmt.Println(testxxx())
	fmt.Println(testxxx() == nil)
	fmt.Println(testxxx() != nil)
}
