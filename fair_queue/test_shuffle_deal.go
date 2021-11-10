package main

import (
	"fmt"
	"math"
	"math/rand"
)

var numQueues uint64

func shuffleDealAndPick(v, nq uint64,
	lengthOfQueue func(int) int,
	mr func( /*in [0, nq-1]*/ int) /*in [0, numQueues-1] and excluding previously determined members of I*/ int,
	nRem, minLen, bestIdx int) int {

	if nRem < 1 {
		return bestIdx
	}
	vNext := v / nq
	ai := int(v - nq*vNext)
	ii := mr(ai)
	i := numQueues - nq // i is used only for debug printing
	mrNext := func(a /*in [0, nq-2]*/ int) /*in [0, numQueues-1] and excluding I[0], I[1], ... ii*/ int {
		if a < ai {
			fmt.Printf("mr[%v](%v) going low\n", i, a)
			return mr(a)
		}
		fmt.Printf("mr[%v](%v) going high\n", i, a)
		return mr(a + 1)
	}
	lenI := lengthOfQueue(ii)
	fmt.Printf("Considering A[%v]=%v, I[%v]=%v, qlen[%v]=%v\n\n", i, ai, i, ii, i, lenI)
	if lenI < minLen {
		minLen = lenI
		bestIdx = ii
	}
	return shuffleDealAndPick(vNext, nq-1, lengthOfQueue, mrNext, nRem-1, minLen, bestIdx)
}

func main() {
	numQueues = uint64(128)
	handSize := 6
	hashValue := rand.Uint64()
	queueIndex := shuffleDealAndPick(hashValue, numQueues, func(idx int) int { return idx % 10 }, func(i int) int { return i }, handSize, math.MaxInt32, -1)
	fmt.Printf("For V=%v, numQueues=%v, handSize=%v, chosen queue is %v\n", hashValue, numQueues, handSize, queueIndex)
}

// Considering A[0]=82, I[0]=82, qlen[0]=2

// mr[0](99) going high
// Considering A[1]=99, I[1]=100, qlen[1]=0

// mr[1](83) going low
// mr[0](83) going high
// Considering A[2]=83, I[2]=84, qlen[2]=4

// mr[2](68) going low
// mr[1](68) going low
// mr[0](68) going low
// Considering A[3]=68, I[3]=68, qlen[3]=8

// mr[3](97) going high
// mr[2](98) going high
// mr[1](99) going high
// mr[0](100) going high
// Considering A[4]=97, I[4]=101, qlen[4]=1

// mr[4](89) going low
// mr[3](89) going high
// mr[2](90) going high
// mr[1](91) going low
// mr[0](91) going high
// Considering A[5]=89, I[5]=92, qlen[5]=2

// For V=5577006791947779410, numQueues=128, handSize=6, chosen queue is 100
