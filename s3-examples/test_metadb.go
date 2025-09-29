package main

import (
	"fmt"
	"log"

	bolt "go.etcd.io/bbolt"
)

func main() {
	db, err := bolt.Open("/var/lib/containerd/io.containerd.metadata.v1.bolt/meta.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	db.View(func(tx *bolt.Tx) error {
		tx.ForEach(func(name []byte, b *bolt.Bucket) error {
			b.ForEachBucket(func(k []byte) error {
				fmt.Printf("bucket=%s\n", name)
				b.ForEachBucket(func(k []byte) error {
					fmt.Printf("--bucket=%s\n", k)
					nestedBucket := b.Bucket(k)
					nestedBucket.ForEachBucket(func(k []byte) error {
						fmt.Printf("----bucket=%s\n", k)
						nestedBucket2 := nestedBucket.Bucket(k)
						nestedBucket2.ForEachBucket(func(k []byte) error {
							fmt.Printf("------bucket=%s\n", k)
							nestedBucket3 := nestedBucket2.Bucket(k)
							nestedBucket3.ForEachBucket(func(k []byte) error {
								fmt.Printf("--------bucket=%s\n", k)
								nestedBucket4 := nestedBucket3.Bucket(k)
								nestedBucket4.ForEach(func(k, v []byte) error {
									fmt.Printf("----------key=%s, value=%s\n", k, v)
									return nil
								})
								return nil
							})
							return nil
						})
						return nil
					})
					return nil
				})
				return nil
			})
			return nil
		})
		return nil
	})
}
