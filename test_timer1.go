package main

import (
	"log"
	"time"
)

func test1() {
	log.Println(time.Now())
	for i := 0; ; i++ {
		timer := time.NewTimer(0)
		stopped := timer.Stop()
		log.Println(stopped)
		var branch string
		select {
		case <-timer.C:
			branch = "recv"
		default:
			branch = "default"
		}
		log.Printf("-----------$$$-----------------------branch = %s--\n", branch)
		log.Println(len(timer.C))
		log.Println(timer.Reset(time.Minute))
		select {
		case <-timer.C:
			log.Println(time.Now())
			log.Fatalf("received from timer channel; i = %d; stopped = %t; branch = %s", i, stopped, branch)
			// default:
		case <-time.After(100 * time.Millisecond):
			log.Println("***after 100 millisecond")
		}
		timer.Stop()
		log.Printf("----------------------------------branch = %s--\n\n", branch)
	}
}

func main() {
	test1()
	test2()

	d := time.Second
	f := func() {}
	time.NewTimer(d)
	time.AfterFunc(d, f)
	time.After(d)
}

func test2() {
	log.Println(time.Now())
	timer := time.NewTimer(time.Second)
	log.Println(timer.Stop())
	timer.Reset(10 * time.Second)
	select {
	case <-timer.C:
		log.Println(time.Now())
	}
	log.Println(time.Now())
}

/*
Example output (go version go1.5 linux/amd64)

$ go run test.go                                                                                              !
2015/09/22 13:22:49 received from timer channel; i = 13674; stopped = false; branch = default
exit status 1
*/
