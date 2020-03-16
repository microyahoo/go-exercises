package main

import (
	"log"
	"time"
)

func main() {
	for i := 0; ; i++ {
		timer := time.NewTimer(0)
		stopped := timer.Stop()
		var branch string
		select {
		case <-timer.C:
			branch = "recv"
		default:
			branch = "default"
		}
		timer.Reset(time.Minute)
		select {
		case <-timer.C:
			log.Fatalf("received from timer channel; i = %d; stopped = %t; branch = %s", i, stopped, branch)
		default:
		}
		timer.Stop()
		// log.Println("----")
	}
}

/*
Example output (go version go1.5 linux/amd64)

$ go run test.go                                                                                              !
2015/09/22 13:22:49 received from timer channel; i = 13674; stopped = false; branch = default
exit status 1
*/
