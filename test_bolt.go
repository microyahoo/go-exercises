package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"

	"github.com/boltdb/bolt"
)

func main() {
	// Remove previous data
	os.Remove("/tmp/test1.db")
	os.Remove("/tmp/test2.db")

	b, err := bolt.Open("/tmp/test1.db", 0600, nil)
	log.Println("Writing data.")
	err = b.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("haha"))
		if err != nil {
			return err
		}
		d := make([]byte, 128)

		for i := 0; i < 10000; i += 1 {
			n, err := rand.Read(d)
			if n != 128 {
				panic("bad len")
			}
			if err != nil {
				return err
			}
			err = b.Put(d, d)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		log.Panic("Inserting.")
	}
	err = b.Close()
	if err != nil {
		panic("can't close file")
	}

	log.Println("Testing")
	data, err := ioutil.ReadFile("/tmp/test1.db")
	if err != nil {
		log.Println(err)
		panic("can't read source db")
	}
	err = ioutil.WriteFile("/tmp/test2.db", data[:len(data)/2], 0600)
	if err != nil {
		panic("can't write source db")
	}
	testDB("/tmp/test2.db")
}

func testDB(fn string) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered in testDB", r)
		}
	}()
	b, err := bolt.Open(fn, 0600, nil)
	if err != nil {
		return
	}
	_ = b
	return
}
