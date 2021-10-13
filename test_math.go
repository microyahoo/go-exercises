package main

import (
	"fmt"
	"math"
	"time"
)

func main() {
	x := 0.210752
	fmt.Println(x * 100)
	fmt.Println(int64(x * 100))
	fmt.Println(math.Round(x * 100))
	fmt.Println("math.Ceil(99.0 / 8) = ", math.Ceil(99.0/8))
	fmt.Println("math.Round(99.0 / 8) = ", math.Round(99.0/8))
	fmt.Printf("%.2f\n", x*100)
	fmt.Printf("%b\n", 4<<10)

	var conns map[string]int
	if n := conns["a"]; n < 1 {
		fmt.Println(n, "xxx")
		// conns["a"] = 1
	}
	fmt.Println(1e9 / 1e8)
	fmt.Println(float64(time.Second / time.Nanosecond))
	fmt.Println(SeatsTimesDuration(12, 10*time.Second))

	ss := SeatsTimesDuration(12, 10*time.Second)
	fmt.Println(ss.ToFloat())
	fmt.Println(ss.DurationPerSeat(12))
	fmt.Println(ss.String())
}

type SeatSeconds uint64

const ssScale = 1e8

func SeatsTimesDuration(seats float64, duration time.Duration) SeatSeconds {
	return SeatSeconds(math.Round(seats * float64(duration/time.Nanosecond) / (1e9 / ssScale)))
}

func (ss SeatSeconds) ToFloat() float64 {
	return float64(ss) / ssScale
}

// DurationPerSeat returns duration per seat.
// This division may lose precision.
func (ss SeatSeconds) DurationPerSeat(seats float64) time.Duration {
	return time.Duration(float64(ss) / seats * (float64(time.Second) / ssScale))
}

func (ss SeatSeconds) String() string {
	const div = SeatSeconds(ssScale)
	quo := ss / div
	rem := ss - quo*div
	return fmt.Sprintf("%d.%08dss", quo, rem)
}
