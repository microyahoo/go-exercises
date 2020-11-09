package main

import (
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

func main() {
	db, err := bolt.Open("/tmp/my.db", 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("MyBucket-1"))

		c := b.Cursor()

		var i int64
		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			i++
			// for k, v := c.First(); k != nil; k, v = c.Next() {
			// fmt.Printf("key=%s, value=%s\n", k, v)
		}

		fmt.Println(i)
		return nil
	})
}
