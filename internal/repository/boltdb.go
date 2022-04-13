package repository

import (
	// "fmt"
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
func (b *BoltDB) Close() error {
	return b.Database.Close()
}
func (b *BoltDB) BucketExist(bucketName string) bool {
	var exist bool
	b.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucketName))
		if b != nil {
			exist = true
		} 
		return nil
	})
	return exist
}
func (b *BoltDB) AddValue(bucketName, key, value []byte) error {
	return b.Database.Batch(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists(bucketName)
		if err != nil {
			return err
		}
		return b.Put(key, value)
	})
}
func (b *BoltDB) GetValue(bucketName, key []byte) (string, error) {
	var v []byte
	err := b.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return nil
		}
		v = b.Get(key)
		return nil
	})
	return string(v), err
}
