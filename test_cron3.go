package main

import (
	"fmt"
	"log"
	"time"

	"github.com/alex023/clock"
)

var (
	myClock = clock.NewClock()
	jobFunc = func() {
		fmt.Println("schedule once")
	}
)

func test1() {
	//add a task that executes once,interval 100 millisecond
	myClock.AddJobWithInterval(time.Duration(100*time.Millisecond), jobFunc)

	//wait a second,watching
	time.Sleep(1 * time.Second)

	//Output:
	//
	//schedule once
}

func ExampleClock_AddJobRepeat() {
	var (
		myClock = clock.NewClock()
	)
	//define a repeat task
	fn := func() {
		fmt.Println("schedule repeat")
	}
	//add in clock,execute three times,interval 200 millisecond
	_, inserted := myClock.AddJobRepeat(time.Duration(time.Millisecond*200), 3, fn)
	if !inserted {
		log.Println("failure")
	}
	//wait a second,watching
	time.Sleep(time.Second)
	//Output:
	//
	//schedule repeat
	//schedule repeat
	//schedule repeat
}

func ExampleClock_RmJob() {
	var (
		myClock = clock.NewClock()
		count   int
		jobFunc = func() {
			count++
			fmt.Println("do ", count)
		}
	)
	//创建任务，间隔1秒，执行两次
	job, _ := myClock.AddJobRepeat(time.Second*1, 2, jobFunc)

	//任务执行前，撤销任务
	time.Sleep(time.Millisecond * 500)
	job.Cancel()

	//等待3秒，正常情况下，事件不会再执行
	time.Sleep(3 * time.Second)

	//Output:
	//
	//
}

func main() {
	test1()
}

// https://github.com/alex023/clock
