package main

import "fmt"

// Greeting function types
// type Greeting func(name string) string

// func say(g Greeting, n string) {
// 	fmt.Println(g(n))
// }

// func english(name string) string {
// 	return "Hello, " + name
// }

// func main() {
// 	say(english, "World")
// }

// =============================================
//A function type denotes the set of all functions with the same parameter and result types.

// Greeting function types
type Greeting func(name string) string

func (g Greeting) say(n string) {
	fmt.Println(g(n))
}

func english(name string) string {
	return "Hello, " + name
}

func main() {
	// g := Greeting(english)
	// g.say("World")
	Greeting(english).say("world")
}
