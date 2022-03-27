package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"sync"
)

type job func(context.Context)

type slice struct {
	pendings []job
	mu       sync.Mutex
}

func (s *slice) add(j job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Println("add job")
	s.pendings = append(s.pendings, j)
}

func (s *slice) finish() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.pendings) != 0 {
		s.pendings = s.pendings[1:]
		fmt.Printf("****finish job: %d\n", len(s.pendings))
	}
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)

	sl := &slice{}

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	go func(s *slice) {
		defer wg.Done()
		for {
			s.add(func(context.Context) {
				fmt.Println("Add job")
			})
		}
	}(sl)

	go func(s *slice) {
		defer wg.Done()
		for {
			s.finish()
		}
	}(sl)

	wg.Wait()
	fmt.Println("End")
}
