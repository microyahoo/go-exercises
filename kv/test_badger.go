package main

import (
	"fmt"
	"unsafe"

	badger "github.com/dgraph-io/badger/v3"
)

func main() {
	var b, c *[4]byte
	a := []byte("1111")
	b = new([4]byte)
	copy(b[:], a)
	c = new([4]byte)
	copy((*c)[:], a)
	fmt.Println(a, b, c)

	//  打开DB
	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	defer db.Close()

	// 读写事务
	err = db.Update(func(txn *badger.Txn) error {
		// Your code here…
		txn.Set([]byte("answer"), []byte("42"))
		txn.Get([]byte("answer"))
		return nil
	})
	// 只读事务
	err = db.View(func(txn *badger.Txn) error {
		// Your code here…
		txn.Get([]byte("answer_v1"))
		return nil
	})
	// 遍历keys
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				fmt.Printf("key=%s, value=%s\n", k, v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	// vlog 的GC
	err = db.RunValueLogGC(0.7)
	_ = err

	fmt.Println("======================================================")
	fmt.Println(int(unsafe.Sizeof(uint32(0))))     // 4
	fmt.Println(int(unsafe.Sizeof(uint64(0))) - 1) // 7
}
