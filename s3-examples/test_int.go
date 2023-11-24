package main

import (
	"fmt"
	"math"
	"time"

	"github.com/mitchellh/hashstructure"
	"k8s.io/apimachinery/pkg/util/wait"
)

func main() {
	a := 1234
	b := 10
	fmt.Println(a / b)

	// The maximum retry duration = initial duration * retry factor ^ # steps. Rearranging, this gives
	// # steps = log(maximum retry / initial duration) / log(retry factor).
	const retryFactor = 1.5
	const initialDurationMs = 100
	maxMs := (5 * time.Second).Milliseconds()
	fmt.Println(maxMs) // 5000
	if maxMs < initialDurationMs {
		maxMs = initialDurationMs
	}
	steps := int(math.Ceil(math.Log(float64(maxMs)/initialDurationMs) / math.Log(retryFactor)))
	if steps < 1 {
		steps = 1
	}
	backoff := wait.Backoff{
		Duration: initialDurationMs * time.Millisecond,
		Factor:   retryFactor,
		Steps:    steps,
	}
	// wait.Backoff{Duration:100000000, Factor:1.5, Jitter:0, Steps:10, Cap:0}
	fmt.Printf("%#v\n", backoff)

	key := map[string]string{"1": "a"}
	hash, _ := hashstructure.Hash(key, nil)
	fmt.Println(hash)
	hash, _ = hashstructure.Hash(map[string]string{"1": "a"}, nil)
	fmt.Println(hash)

	// s := "abcdefg"
	// fmt.Println(s[-1:])

	tmp := []string{"a"}
	fmt.Println(append(tmp[:0], tmp[1:]...))
}
