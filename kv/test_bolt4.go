package main

import (
	"encoding/binary"
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
		// for k, _ := c.Seek(itob(uint64(3))); k != nil; k, _ = c.Next() {
		for k, _ := c.Seek([]byte("xxxxx")); k != nil; k, _ = c.Next() {
			// for k, _ := c.First(); k != nil; k, _ = c.Next() {
			i++
			// for k, v := c.First(); k != nil; k, v = c.Next() {
			// fmt.Printf("key=%s, value=%s\n", k, v)
			fmt.Printf("key=%d\n", btoi(k))
		}

		fmt.Println(i)
		return nil
	})
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

// btoi returns an 8-byte big endian representation of v.
func btoi(b []byte) uint64 {
	var i uint64
	i = binary.BigEndian.Uint64(b)
	return i
}
