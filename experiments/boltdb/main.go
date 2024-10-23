package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"log"
	"os"
)

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	defer os.Remove("my.db")

	err = db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("MyBucket"))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}

		err = b.Put([]byte("answer"), []byte("42"))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	err = db.View(func(tx *bolt.Tx) error {
		answer := tx.Bucket([]byte("MyBucket")).Get([]byte("answer"))
		fmt.Printf("The answer is: %s\n", answer)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}
