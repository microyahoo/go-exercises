package main

import (
	"fmt"
	"time"
)

type DB struct {
	num int
}

type Tx struct {
	writable bool
	pending  int
	db       *DB
}

func main() {
	tx := &Tx{
		writable: true,
		pending:  10,
		db:       &DB{},
	}
	go func(tx *Tx) {
		time.Sleep(1 * time.Second)
		fmt.Printf("tx: %p, &tx: %p,  %v, %v, %v\n", tx, &tx, tx, tx.pending, tx.db)
	}(tx)
	fmt.Printf("before: %p\n", tx)
	fmt.Printf("before&: %p\n", &tx)
	tx = nil
	fmt.Printf("after: %p\n", tx)

	fmt.Println(100)
	time.Sleep(100 * time.Second)
}
