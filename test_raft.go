package main

import (
	"fmt"
	"time"
)

type message struct {
	Type string
	Term int64
	From int64
	To   int64
}

type raft struct {
	msgs     []message
	readyc   chan ready
	advancec chan struct{}
}

func (r *raft) acceptReady() {
	r.msgs = nil
}

func (r *raft) hasReady() bool {
	return len(r.msgs) > 0
}

func (r *raft) newReady() ready {
	return ready{
		msgs: r.msgs,
	}
}

func (r *raft) advance() {
	r.advancec <- struct{}{}
}

func (r *raft) ready() <-chan ready {
	return r.readyc
}

type ready struct {
	msgs []message
}

func main() {
	var raft = &raft{
		readyc:   make(chan ready),
		advancec: make(chan struct{}),
		msgs: []message{
			{
				Type: "normal",
				Term: 1,
				From: 1,
				To:   2,
			},
			{
				Type: "confChange",
				Term: 1,
				From: 3,
				To:   2,
			},
			{
				Type: "snapshot",
				Term: 2,
				From: 5,
				To:   2,
			},
		},
	}

	go func() {
		for {
			select {
			case rd := <-raft.ready():
				time.Sleep(time.Second * 5)
				fmt.Printf("After 5s, ready = %#v\n", rd)
				raft.advance()
			}
		}
	}()

	go func() {
		var readyc chan ready
		var advancec chan struct{}
		var rd ready

		for {
			if advancec != nil {
				readyc = nil
			} else if raft.hasReady() {
				rd = raft.newReady()
				readyc = raft.readyc
			}
			select {
			case readyc <- rd:
				raft.acceptReady()
				advancec = raft.advancec
				fmt.Println("After acceptReady")
			case <-advancec:
				// Advance
				rd = ready{}
				advancec = nil
				fmt.Println("Advance to next run")
			}
		}
	}()

	select {}
}
