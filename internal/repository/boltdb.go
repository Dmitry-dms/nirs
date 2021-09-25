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

func (b *BoltDB) AddValue(key, value []byte) error {
	b.Database.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(b.Names))
		if err != nil {
			return err
		}
		return b.Put(key, value)
	})
	return nil
}
func (b *BoltDB) GetValue(key []byte) (string, error) {
	var v []byte
	err := b.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(b.Names))
		v = b.Get(key)
		return nil
	})
	return string(v), err
}
