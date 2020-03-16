package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func test1() {
	timer := time.NewTimer(time.Second * 2)
	<-timer.C
	fmt.Println("Timer expired")
}

func test2() {
	timer := time.NewTimer(time.Second)
	go func() {
		<-timer.C
		fmt.Println("Timer expired")
	}()
	stop := timer.Stop()
	fmt.Println("Timer cancelled:", stop)
}

func test3() {
	ticker := time.NewTicker(time.Millisecond * 500)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()
	time.Sleep(time.Millisecond * 1500)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}

func test4() {
	timeChan := time.NewTimer(time.Second).C

	tickChan := time.NewTicker(time.Millisecond * 400).C

	doneChan := make(chan bool)
	go func() {
		time.Sleep(time.Second * 2)
		doneChan <- true
	}()

	for {
		select {
		case <-timeChan:
			fmt.Println("Timer expired")
		case <-tickChan:
			fmt.Println("Ticker ticked")
		case <-doneChan:
			fmt.Println("Done")
			return
		}
	}
}

func test5() {
	now := time.Now()
	after := now.AddDate(0, 0, 825)
	before := now.Add(-5 * time.Second)
	duration, _ := time.ParseDuration("1h")
	fmt.Println(duration)
	fmt.Println(duration.Seconds())
	fmt.Println("-----------1--------")
	fmt.Println(now)
	fmt.Println(before)
	fmt.Println(after)
	duration = before.Sub(now)
	fmt.Println(duration)
	timer := time.NewTimer(duration)
	for {
		select {
		case now = <-timer.C:
			fmt.Printf("*****timer %v\n", now)
			// default:
			// 	fmt.Println("hello")
		}
		break
	}
}

func keysHelper(keys ...string) string {
	return strings.Join(keys, "-")
}

type Job struct {
	a string
}

func (job *Job) init(a string) {
	job.a = a
}

func main() {
	test5()
	var job *Job
	job = new(Job)
	job.init("xxx")
	fmt.Println(job)
	fmt.Println(keysHelper("a", "b"))
	fmt.Println(keysHelper("b", strconv.FormatInt(100, 10)))
	fmt.Println(keysHelper("a"))
	test1()
	test2()
	test3()
	test4()
}

// https://mmcgrana.github.io/2012/09/go-by-example-timers-and-tickers.html
