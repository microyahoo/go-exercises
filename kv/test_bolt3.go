package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	bolt "go.etcd.io/bbolt"
)

var (
	dbPath     = "test.db"
	bucketName = "test"

	numKeys = 100

	keyLen = 100
	valLen = 750

	keys = make([][]byte, numKeys)
	vals = make([][]byte, numKeys)
)

func init() {
	fmt.Println("Starting generating random data...")
	for i := range keys {
		keys[i] = randBytes(keyLen)
		vals[i] = randBytes(valLen)
	}
	fmt.Println("Done")
}

func main() {
	defer os.Remove(dbPath)

	initMmapSize := 1 << 31 // 2GB
	initMmapSize = 0        // comment this out to not block write

	var boltOpenOptions = &bolt.Options{InitialMmapSize: initMmapSize}

	db, err := bolt.Open(dbPath, 0600, boltOpenOptions)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Println("Creating bucket:", bucketName)
	if err := db.Update(func(tx *bolt.Tx) error {
		if _, err := tx.CreateBucket([]byte(bucketName)); err != nil {
			return err
		}
		return nil
	}); err != nil {
		panic(err)
	}
	fmt.Println("Done with creating bucket:", bucketName)

	fmt.Println("Starting reading...")
	go func() {
		_, err := db.Begin(false)
		if err != nil {
			panic(err)
		}
		fmt.Println("before select")
		select {}
		fmt.Println("after select")
	}()

	fmt.Println("Starting writing...")
	for {
		for j := range keys {
			tw := time.Now()
			if err := db.Update(func(tx *bolt.Tx) error {
				b := tx.Bucket([]byte(bucketName))
				if err := b.Put(keys[j], vals[j]); err != nil {
					return err
				}
				return nil
			}); err != nil {
				panic(err)
			}
			fmt.Printf("#%d write took: %v\n", j, time.Since(tw))
		}
	}
}

func randBytes(n int) []byte {
	const (
		letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
		letterIdxBits = 6                    // 6 bits to represent a letter index
		letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
		letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
	)
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return b
}
