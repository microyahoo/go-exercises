package main

import (
	"fmt"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

func main() {
	db, _ := leveldb.OpenFile("db", nil)
	defer db.Close()
	//读写数据库:
	_ = db.Put([]byte("key1"), []byte("好好检查"), nil)
	_ = db.Put([]byte("key2"), []byte("天天向上"), nil)
	_ = db.Put([]byte("key:3"), []byte("就会一个本事"), nil)
	data, _ := db.Get([]byte("key1"), nil)
	fmt.Println(string(data))
	//迭代数据库内容:
	iter := db.NewIterator(nil, nil)
	for iter.Next() {
		key := iter.Key()
		value := iter.Value()
		fmt.Println(string(key), string(value))
	}
	iter.Release()
	iter.Error()
	//Seek-then-Iterate:
	iter = db.NewIterator(nil, nil)
	for ok := iter.Seek([]byte("key:")); ok; ok = iter.Next() {
		// Use key/value.
		fmt.Println("Seek-then-Iterate:")
		fmt.Println(string(iter.Value()))
	}
	iter.Release()
	//Iterate over subset of database content:
	iter = db.NewIterator(&util.Range{Start: []byte("key:"), Limit: []byte("xoo")}, nil)
	for iter.Next() {
		// Use key/value.
		fmt.Println("Iterate over subset of database content:")
		fmt.Println(string(iter.Value()))
	}
	iter.Release()
	//Iterate over subset of database content with a particular prefix:
	iter = db.NewIterator(util.BytesPrefix([]byte("key")), nil)
	for iter.Next() {
		// Use key/value.
		fmt.Println("Iterate over subset of database content with a particular prefix:")
		fmt.Println(string(iter.Value()))
	}
	iter.Release()
	_ = iter.Error()
	//批量写:
	batch := new(leveldb.Batch)
	batch.Put([]byte("foo"), []byte("value"))
	batch.Put([]byte("bar"), []byte("another value"))
	batch.Delete([]byte("baz"))
	_ = db.Write(batch, nil)
	_ = db.Delete([]byte("key"), nil)
}
