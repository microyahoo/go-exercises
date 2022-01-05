package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func main() {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:12379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	lockKey := "/lock"

	session, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}

	go func(session *concurrency.Session) {
		m := concurrency.NewMutex(session, lockKey)
		if err := m.Lock(context.TODO()); err != nil {
			log.Fatal("go1 get mutex failed " + err.Error())
		}
		fmt.Printf("go1 get mutex sucess\n")
		fmt.Println(m)
		time.Sleep(time.Duration(10) * time.Second)
		m.Unlock(context.TODO())
		fmt.Printf("go1 release lock\n")
	}(session)

	go func() {
		time.Sleep(time.Duration(2) * time.Second)
		session, err := concurrency.NewSession(cli)
		if err != nil {
			log.Fatal(err)
		}
		m := concurrency.NewMutex(session, lockKey)
		if err := m.Lock(context.TODO()); err != nil {
			log.Fatal("go2 get mutex failed " + err.Error())
		}
		fmt.Printf("go2 get mutex sucess\n")
		fmt.Println(m)
		time.Sleep(time.Duration(2) * time.Second)
		m.Unlock(context.TODO())
		fmt.Printf("go2 release lock\n")
	}()

	// 可重入锁
	go func(session *concurrency.Session) {
		m := concurrency.NewMutex(session, lockKey)
		if err := m.Lock(context.TODO()); err != nil {
			log.Fatal("go3 get mutex failed " + err.Error())
		}
		fmt.Printf("go3 get mutex sucess\n")
		fmt.Println(m)
		time.Sleep(time.Duration(5) * time.Second)
		fmt.Printf("go3 exit lock\n")
	}(session)

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
