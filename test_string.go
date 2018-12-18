package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	fmt.Println(strings.Split("network__addresses", "__"))
	fmt.Println(0x7FFFFFF)

	rot13 := func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			r += 'a' - 'A'
			return r
			// return 'A' + (r-'A'+13)%26
		case r >= 'a' && r <= 'z':
			r -= 'a' - 'A'
			return r
			// return 'a' + (r-'a'+13)%26
		}
		return r
	}
	fmt.Println(strings.Map(rot13, "'Twas brillig and the slithy gopher..."))

	GuessingGame()
}

func GuessingGame() {
	var s string
	fmt.Printf("Pick an integer from 0 to 100.\n")
	answer := sort.Search(100, func(i int) bool {
		fmt.Printf("Is your number <= %d? ", i)
		fmt.Scanf("%s", &s)
		fmt.Println(s)
		return s != "" && s[0] == 'y'
	})
	fmt.Printf("Your number is %d.\n", answer)
}
