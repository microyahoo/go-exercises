package main

import "fmt"

type Greeting struct {
	say func(name string) string
}

func newGreeting(f func(string) string) *Greeting {
	return &Greeting{say: f}
}

func (g *Greeting) exclamation(name string) string { return g.say(name) + "!" }

func main() {
	english := &Greeting{say: func(name string) string {
		return "Hello, " + name
	}}

	french := newGreeting(func(name string) string {
		return "Bonjour, " + name
	})

	fmt.Println(english.exclamation("ANisus"))
	fmt.Println(french.exclamation("ANisus"))

	var x *int
	var y *int
	fmt.Println(x == y)
}
