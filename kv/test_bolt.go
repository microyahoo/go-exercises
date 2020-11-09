package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	// "os"
	"strconv"
	// "time"
	"unsafe"

	bolt "go.etcd.io/bbolt"
)

func main() {
	var ids []uint64
	fmt.Println(unsafe.Sizeof(ids[0]))
	fmt.Println(unsafe.Sizeof(ids[1]))
	// fmt.Println(ids[0])

	db, err := bolt.Open("/tmp/my.db", 0600, nil)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	// go func() {
	// 	// Grab the initial stats.
	// 	prev := db.Stats()

	// 	for {
	// 		// Wait for 10s.
	// 		time.Sleep(10 * time.Second)

	// 		// Grab the current stats and diff them.
	// 		stats := db.Stats()
	// 		diff := stats.Sub(&prev)

	// 		// Encode stats to JSON and print to STDERR.
	// 		json.NewEncoder(os.Stderr).Encode(diff)

	// 		// Save stats for the next loop.
	// 		prev = stats
	// 	}
	// }()

	db.Update(func(tx *bolt.Tx) error {
		for i := 0; i < 10; i++ {
			b, err := tx.CreateBucketIfNotExists([]byte(fmt.Sprintf("%s-%d", "MyBucket", i)))
			if err != nil {
				log.Panic(fmt.Errorf("create bucket: %s", err))
			}
			b.Put([]byte("my-key"), []byte("my-value"))
			// fmt.Println(b.Put([]byte("my-key"), []byte("my-value")))
			subBucket, err := b.CreateBucketIfNotExists([]byte("sub-bucket"))
			if err != nil {
				log.Panic(fmt.Errorf("create sub bucket: %s", err))
			}
			err = subBucket.Put([]byte("key-hello"), []byte("value-hello"))
			if err != nil {
				log.Panic(fmt.Errorf("put key/value to sub bucket: %s", err))
			}
			err = subBucket.Put([]byte("key-hello2"), []byte("value-hello2"))
			if err != nil {
				log.Panic(fmt.Errorf("put key/value to sub bucket: %s", err))
			}

			k := make([]byte, 12)
			v := make([]byte, 5120)

			for i := 0; i < 10; i++ {
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
		}
		return nil
	})

	// db.View(func(tx *bolt.Tx) error {
	// 	b := tx.Bucket([]byte("MyBucket"))
	// 	v := b.Get([]byte("my-key"))
	// 	fmt.Printf("The answer is: %s\n", v)
	// 	return nil
	// })

	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			log.Panic(fmt.Errorf("create bucket: %s", err))
		}
		// b.Put([]byte("my-key"), []byte("my-value"))
		return nil
	})

	store := &Store{
		db: db,
	}

	user1 := &User{
		Name: "zhangsan",
	}
	store.CreateUser(user1)
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		v := b.Get(itob(user1.ID))
		fmt.Printf("The answer is: %s\n", v)
		return nil
	})

	// fmt.Println("-----------------iterate MyBucket---------------------")
	// db.View(func(tx *bolt.Tx) error {
	// 	// Assume bucket exists and has keys
	// 	b := tx.Bucket([]byte("MyBucket"))
	// 	c := b.Cursor()
	// 	for k, v := c.First(); k != nil; k, v = c.Next() {
	// 		fmt.Printf("key=%s, value=%s\n", k, v)
	// 	}
	// 	return nil
	// })

	// fmt.Println("-----------------iterate users---------------------")
	// db.View(func(tx *bolt.Tx) error {
	// 	// Assume bucket exists and has keys
	// 	b := tx.Bucket([]byte("users"))
	// 	c := b.Cursor()
	// 	for k, v := c.First(); k != nil; k, v = c.Next() {
	// 		fmt.Printf("key=%v, value=%s\n", btoi(k), v)
	// 	}
	// 	return nil
	// })

	// select {}

}

type Store struct {
	db *bolt.DB
}

// CreateUser saves u to the store. The new user ID is set on u once the data is persisted.
func (s *Store) CreateUser(u *User) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		// Retrieve the users bucket.
		// This should be created when the DB is first opened.
		b := tx.Bucket([]byte("users"))

		// Generate ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call
		id, err := b.NextSequence()
		if err != nil {
			return err
		}
		u.ID = id

		// Marshal user data into bytes.
		buf, err := json.Marshal(u)
		if err != nil {
			return err
		}

		// Persist bytes to users bucket.
		return b.Put(itob(u.ID), buf)
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

type User struct {
	ID   uint64
	Name string
}

// createUser creates a new user in the given account.
func (s *Store) CreateUserByAccountID(accountID uint64, u *User) error {
	// Start the transaction.
	tx, err := s.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Retrieve the root bucket for the account.
	// Assume this has already been created when the account was set up.
	root := tx.Bucket([]byte(strconv.FormatUint(accountID, 10)))

	// Setup the users bucket.
	bkt, err := root.CreateBucketIfNotExists([]byte("USERS"))
	if err != nil {
		return err
	}

	// Generate an ID for the new user.
	userID, err := bkt.NextSequence()
	if err != nil {
		return err
	}
	u.ID = userID

	// Marshal and save the encoded user.
	if buf, err := json.Marshal(u); err != nil {
		return err
	} else if err := bkt.Put([]byte(strconv.FormatUint(u.ID, 10)), buf); err != nil {
		return err
	}

	// Commit the transaction.
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}
