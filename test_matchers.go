package main

import (
	"fmt"
)

type Matchers []*Matcher

type Ints []int

type Matcher struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Ints  Ints
}

type LabelName string
type LabelValue string

type LabelSet map[LabelName]LabelValue

func (m *Matcher) test() bool {
	return true
}

func (ms Matchers) Match(lset LabelSet) bool {
	return true
}

func (ms *Matchers) MatchP(lset LabelSet) bool {
	return true
}
func main() {
	var matchers Matchers
	fmt.Println(matchers)
	fmt.Println(len(matchers))
	fmt.Println(matchers == nil)
	fmt.Println(matchers.Match(LabelSet(nil)))
	fmt.Printf("%p, %p\n", nil, matchers)
	var matchersP *Matchers
	fmt.Println(matchersP)
	fmt.Println(matchersP == nil)
	fmt.Println(matchersP.MatchP(LabelSet(nil)))

	var matcher *Matcher
	fmt.Println(matcher)
	// fmt.Println(matcher == nil)
	if matcher == nil {
		fmt.Println("matcher is nil")
	}
	fmt.Println(matcher.test())
}
