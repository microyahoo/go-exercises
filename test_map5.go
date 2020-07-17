package main

import (
	log "github.com/Sirupsen/logrus"
	"math/rand"
	"os"
	"os/signal"
	"time"
)

func main() {
	log.Info("=== START ===")
	defer func() { log.Info("=== DONE ===") }()

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

	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt)
	for {
		select {
		case <-osSignals:
			log.Info("Recieved ^C command. Exit")
			return
		}
	}
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
