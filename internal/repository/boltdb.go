package repository

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
)

type BoltDB struct {
	FileName                   string
	Database                   *bolt.DB
	Names, Addresses, Passport string
}

func NewBoltDB(fileName string) *BoltDB {
	db, err := bolt.Open(fileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	return &BoltDB{
		Database:  db,
		Names:     "names",
		Addresses: "address",
		Passport:  "passport"}
}

func (b *BoltDB) AddValue(key, value string) error {
	b.Database.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(b.Names))
		if err != nil {
			return err
		}
		return b.Put([]byte(key), []byte(value))
	})
	return nil
}
func (b *BoltDB) GetValue(key string) (string, error) {
	var v []byte
	err := b.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.Names))
		v = b.Get([]byte(key))
		return nil
	})
	return string(v), err
}
