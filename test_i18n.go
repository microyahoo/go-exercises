package main

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func main() {
	p := message.NewPrinter(language.BritishEnglish)
	p.Printf("There are %v flowers in our garden.\n", 1500)
	p = message.NewPrinter(language.Greek)
	p.Printf("There are %v flowers in our garden.", 1500)
}
