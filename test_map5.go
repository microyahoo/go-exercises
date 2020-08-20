package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.Println("=== START ===")
	defer func() { log.Println("=== DONE ===") }()

	go func() {
		m := make(map[string]string)
		for {
			k := generateRandStr(1024)
			m[k] = generateRandStr(1024 * 1024)

			for k2 := range m {
				delete(m, k2)
				break
			}
		}
	}()

	var y *xx
	x := getNil()
	// y := (*xx)(x)
	if x != nil {
		fmt.Println("=====")
		y = x.(*xx)
	}
	fmt.Println(y)
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt)
	for {
		select {
		case <-osSignals:
			log.Println("Recieved ^C command. Exit")
			return
		}
	}
}

type xx struct {
	a, b string
}

func getNil() interface{} {
	return nil
}

func generateRandStr(n int) string {
	rand.Seed(time.Now().UnixNano())
	const letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
