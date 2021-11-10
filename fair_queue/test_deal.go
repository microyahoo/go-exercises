// https://github.com/kubernetes/enhancements/tree/master/keps/sig-api-machinery/1040-priority-and-fairness

package main

import (
	"fmt"
	"sort"
)

// sum computes the sum of the given slice of numbers
func sum(v []float64) float64 {
	c := append([]float64{}, v...)
	sort.Float64s(c) // to minimize loss of accuracy when summing
	var s float64
	for i := 0; i < len(c); i++ {
		s += c[i]
	}
	return s
}

// choose returns the number of subsets of size m of a set of size n
func choose(n, m int) float64 {
	if m == 0 || m == n {
		return 1
	}
	var ans = float64(n)
	for i := 1; i < m; i++ {
		ans = ans * float64(n-i) / float64(i+1)
	}
	return ans
}

// nthDeal analyzes the result of another shuffle and deal in a series of shuffles and deals.
// Each shuffle and deal randomly picks `handSize` distinct cards from a deck of size `deckSize`.
// Each successive shuffle and deal is independent of previous deals.
// `first` indicates that this is the first shuffle and deal.
// `prevDist[nUnique]` is the probability that the number of unique cards previously dealt is `nUnique`,
// and is unused when `first`.
// `dist[nUnique]` is the probability that the number of unique cards dealt up through this deal is `nUnique`.
// `distSum` is the sum of `dist`, and should be 1.
// `expectedUniques` is the expected value of nUniques at the end of this deal.
// `probNextCovered` is the probability that another shuffle and deal will deal only cards that have already been dealt.
func nthDeal(first bool, handSize, deckSize int, prevDist []float64) (dist []float64, distSum, expectedUniques, probNextCovered float64) {
	dist = make([]float64, deckSize+1)
	expects := make([]float64, deckSize+1)
	nexts := make([]float64, deckSize+1)
	if first {
		dist[handSize] = 1
		expects[handSize] = float64(handSize)
		nexts[handSize] = 1 / choose(deckSize, handSize)
	} else {
		for nUnique := handSize; nUnique <= deckSize; nUnique++ {
			conts := make([]float64, handSize+1)
			for news := 0; news <= handSize; news++ {
				// one way to get to nUnique is for `news` new uniques to appear in this deal,
				// and all the previous deals to have dealt nUnique-news unique cards.
				prevUnique := nUnique - news
				ways := choose(deckSize-prevUnique, news) * choose(prevUnique, handSize-news)
				conts[news] = ways * prevDist[prevUnique]
				//fmt.Printf("nUnique=%v, news=%v, ways=%v\n", nUnique, news, ways)
			}
			dist[nUnique] = sum(conts) / choose(deckSize, handSize)
			expects[nUnique] = dist[nUnique] * float64(nUnique)
			nexts[nUnique] = dist[nUnique] * choose(nUnique, handSize) / choose(deckSize, handSize)
		}

	}
	return dist, sum(dist), sum(expects), sum(nexts)
}

func main() {
	handSize := 7
	deckSize := 256
	fmt.Printf("choose(%v, %v) = %v\n", deckSize, handSize, choose(deckSize, handSize))
	var dist []float64
	var probNextCovered float64
	for nHands := 1; probNextCovered < 0.01; nHands++ {
		var distSum, expected float64
		dist, distSum, expected, probNextCovered = nthDeal(nHands == 1, handSize, deckSize, dist)
		fmt.Printf("After %v hands, distSum=%v, expected=%v, probNextCovered=%v, dist=%v\n", nHands, distSum, expected, probNextCovered, dist)
	}
}

// choose(256, 7) = 1.3161885792e+13
// After 1 hands, distSum=1, expected=7, probNextCovered=7.597695465552631e-14, dist=[0 0 0 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
// After 2 hands, distSum=1, expected=13.808593750000002, probNextCovered=2.365999992609045e-10, dist=[0 0 0 0 0 0 0 7.597695465552631e-14 1.3242783196458237e-10 4.926315349082464e-08 6.759999395685381e-06 0.0004157399628346509 0.012222754907338738 0.16568623318836956 0.8216684625464041 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
// After 3 hands, distSum=1, expected=20.431015014648438, probNextCovered=7.2764682268046046e-09, dist=[0 0 0 0 0 0 0 5.7724976387279015e-27 9.055317045872458e-23 2.0835278376214098e-19 1.570601544390448e-16 5.3084534585810687e-14 9.468801451127312e-12 9.812229814502437e-10 6.26967443461761e-08 2.5610950557747708e-06 6.813158998015206e-05 0.0011831308968703333 0.013193660417613478 0.09036581264178295 0.3435833838946934 0.5516032557765146 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
