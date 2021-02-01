package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"sync"
	"time"
)

func main() {
	// https://golang.org/doc/diagnostics.html
	runtime.SetBlockProfileRate(1)
	runtime.SetMutexProfileFraction(1)
	go func() {
		ip := "0.0.0.0:6060"
		if err := http.ListenAndServe(ip, nil); err != nil {
			fmt.Printf("start pprof failed on %s\n", ip)
			os.Exit(1)
		}
	}()

	_ = http.FS

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	go testWaitGroup()

	for {
		select {
		case <-ticker.C:
			fmt.Printf("goroutines: %d\n", runtime.NumGoroutine())
			// runtime.ReadMemStats()
		}
	}

	// for range ticker.C {
	// 	fmt.Println("hello")
	// }
}

func testWaitGroup() {
	wg := sync.WaitGroup{}

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go calc(&wg, i)
	}

	wg.Wait()
	fmt.Println("all goroutine finish")
}

func calc(w *sync.WaitGroup, i int) {
	fmt.Println("calc: ", i)
	time.Sleep(time.Second * 2)
	w.Done()
}
