package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func play(p *player, table chan int) {
	for {
		ball, ok := <-table
		if !ok {
			fmt.Printf("%s wins.\n", p.name)
			return
		}
		r := rand.Intn(100)
		if r > p.successRatio {
			fmt.Printf("%s loses.\n", p.name)
			close(table)
			return
		}
		fmt.Printf("%d %s\n", ball, p.name)
		time.Sleep(time.Second)
		ball++
		table <- ball
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

type player struct {
	name         string
	successRatio int
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	table := make(chan int)
	ball := 1

	go func() {
		play(&player{name: "Zhang", successRatio: 80}, table)
		wg.Done()
	}()

	go func() {
		time.Sleep(time.Second)
		play(&player{name: "Li", successRatio: 90}, table)
		wg.Done()
	}()

	table <- ball

	wg.Wait()
	fmt.Println("Game over!")
}
