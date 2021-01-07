package main

import (
	"fmt"
	"time"
)

func Test1() {
	var waitFiveHundredMillisections int64 = 500

	startingTime := time.Now().UTC()
	time.Sleep(100 * time.Millisecond)
	endingTime := time.Now().UTC()

	var duration time.Duration = endingTime.Sub(startingTime)
	var durationAsInt64 = int64(duration)

	if durationAsInt64 >= waitFiveHundredMillisections {
		fmt.Printf("Time Elapsed : Wait[%d] Duration[%d]\n", waitFiveHundredMillisections, durationAsInt64)
	} else {
		fmt.Printf("Time DID NOT Elapsed : Wait[%d] Duration[%d]\n", waitFiveHundredMillisections, durationAsInt64)
	}
}

func Test2() {
	var duration_Milliseconds time.Duration = 500 * time.Millisecond
	var duration_Seconds time.Duration = (1250 * 10) * time.Millisecond
	var duration_Minute time.Duration = 2 * time.Minute

	fmt.Printf("Milli [%v]\nSeconds [%v]\nMinute [%v]\n", duration_Milliseconds, duration_Seconds, duration_Minute)
}

func Test3() {
	var duration_Seconds time.Duration = (1250 * 10) * time.Millisecond
	var duration_Minute time.Duration = 2 * time.Minute

	var float64_Seconds float64 = duration_Seconds.Seconds()
	var float64_Minutes float64 = duration_Minute.Minutes()

	fmt.Printf("Seconds [%.3f]\nMinutes [%.2f]\n", float64_Seconds, float64_Minutes)
}

func Test4() {
	var waitFiveHundredMillisections time.Duration = 500 * time.Millisecond

	startingTime := time.Now().UTC()
	time.Sleep(600 * time.Millisecond)
	endingTime := time.Now().UTC()

	var duration time.Duration = endingTime.Sub(startingTime)

	if duration >= waitFiveHundredMillisections {
		fmt.Printf("Wait %v\nNative [%v]\nMilliseconds [%d]\nSeconds [%.3f]\n", waitFiveHundredMillisections, duration, duration.Nanoseconds()/1e6, duration.Seconds())
	}
}

func main() {
	fmt.Println(time.Now().UTC().Format("20190222T114211"))
	fmt.Println(time.Now().UTC().Format(time.UnixDate))
	fmt.Println(time.Now().UTC().Format(time.RFC3339))
	fmt.Println("------++--------")
	fmt.Println(time.Now().Format(time.RFC3339))
	fmt.Println("------++--------")
	fmt.Println(time.Now().UTC().Format(time.RFC3339Nano))
	fmt.Println(time.Now())
	fmt.Println(time.Now().UTC())
	fmt.Println("--------------")
	Test1()
	Test2()
	Test3()
	Test4()
}
