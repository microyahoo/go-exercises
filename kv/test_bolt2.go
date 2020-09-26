package main

import (
	"io/ioutil"
	"log"
	"math/rand"
	"os"

	bolt "go.etcd.io/bbolt"
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
		k := make([]byte, 12)
		v := make([]byte, 5120)

		for i := 0; i < 1000; i += 1 {
			n, err := rand.Read(k)
			if n != 12 {
				panic("bad len")
			}
			if err != nil {
				return err
			}
			n, err = rand.Read(v)
			if n != 5120 {
				panic("bad len")
			}
			if err != nil {
				return err
			}
			err = b.Put(k, v)
			if err != nil {
				return err
			}
		}
		k = make([]byte, 2)
		v = make([]byte, 20)
		for i := 0; i < 1000; i += 1 {
			n, err := rand.Read(k)
			if n != 2 {
				panic("bad len")
			}
			if err != nil {
				return err
			}
			n, err = rand.Read(v)
			if n != 20 {
				panic("bad len")
			}
			if err != nil {
				return err
			}
			err = b.Put(k, v)
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
