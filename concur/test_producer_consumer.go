package main

import (
	"fmt"
	"sync"
	"time"
)

const (
	NumOfProducer = 4
	NumOfConsumer = 5
)

type Producer struct {
	producerId int
}

type Consumer struct {
	consumerId int
}

func newProducer(pid int) *Producer {
	return &Producer{producerId: pid}
}

func newConsumer(cid int) *Consumer {
	return &Consumer{consumerId: cid}
}

func (p *Producer) Run(ch chan<- string) {
	for i := 0; i < 5; i++ {
		fmt.Printf("producer %d put data %d\n", p.producerId, i)
		data := fmt.Sprintf("data %d from Producer %d", i, p.producerId)
		ch <- data
		time.Sleep(time.Millisecond * 1)
	}
}

func (c *Consumer) Run(ch <-chan string) {
	for data := range ch {
		fmt.Printf("Consumer %d got %s\n", c.consumerId, data)
		time.Sleep(time.Millisecond * 1000)
	}
	fmt.Printf("Consumer %d: detect channel closed\n", c.consumerId)
	// for {
	// 	data, ok := <-ch
	// 	if !ok {
	// 		fmt.Printf("Consumer %d: detect channel closed\n", c.consumerId)
	// 		return
	// 	}
	// 	fmt.Printf("Consumer %d got %s\n", c.consumerId, data)
	// 	time.Sleep(time.Millisecond * 1000)
	// }
}

func main() {
	buffer := make(chan string, 10)

	prodWg := sync.WaitGroup{}
	prodWg.Add(NumOfProducer)

	for i := 0; i < NumOfProducer; i++ {
		go func(id int) {
			p := newProducer(id)
			p.Run(buffer)
			prodWg.Done()
		}(i + 1)
	}

	consWg := sync.WaitGroup{}
	consWg.Add(NumOfConsumer)
	for i := 0; i < NumOfConsumer; i++ {
		go func(id int) {
			p := newConsumer(id)
			p.Run(buffer)
			consWg.Done()
		}(i + 1)
	}

	prodWg.Wait()
	close(buffer)

	consWg.Wait()
	fmt.Printf("exit....\n")
}
