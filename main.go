package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dgraph-io/badger"
)

var BadgerDB *badger.DB

func InitBadger() {

	db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))

	if err != nil {
		log.Fatal(err)
	}

	BadgerDB = db

	txn := db.NewTransaction(true)

	testbytes, _ := json.Marshal(map[string]interface{}{
		"uae": "1221",
	})

	err = txn.Set([]byte("test3"), testbytes)

	if err != nil {
		fmt.Println(err)
	}

	txn.Commit()
	testvalue, err := Get("test3")
	if err != nil {
		fmt.Println(err, "get")
	}
	var getinterface interface{}
	_ = json.Unmarshal(testvalue, &getinterface)
	fmt.Println(testvalue)
	fmt.Println(getinterface)

}

func main() {

	opts := badger.DefaultOptions("./tmp/badger")
	opts.Truncate = true
	opts.SyncWrites = true
	opts.NumVersionsToKeep = 1
	opts.CompactL0OnClose = true
	opts.NumLevelZeroTables = 1
	opts.NumLevelZeroTablesStall = 2
	opts.ValueLogFileSize = 1024 * 1024 * 50

	db, err := badger.Open(opts)

	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		// err = txn.Set([]byte("newkey3"), []byte("132132213"))
		// Handle(err)

		// return err
		item, err := txn.Get([]byte("newkey3"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			var getinterface interface{}
			_ = json.Unmarshal(val, &getinterface)
			fmt.Println(val)
			fmt.Println(getinterface)
			return err
		})
		return err
	})
	fmt.Println(err)
}

func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func Get(key string) (values []byte, errors error) {

	err := BadgerDB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			values = make([]byte, len(val))
			copy(values, val)
			return nil
		})

		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

func Set(key string, value map[string]interface{}) error {

	_, err := json.Marshal(value)

	if err != nil {
		return err
	}

	txn := BadgerDB.NewTransaction(true)
	defer txn.Discard()

	err = txn.Set([]byte("test"), []byte("5"))

	if err != nil {
		return err
	}

	txn.Commit()

	return nil
}
