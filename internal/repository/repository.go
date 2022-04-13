package repository

import "database/sql"

type KVRepository interface {
	AddValue(bucketName, key, value []byte) error
	GetValue(bucketName, key []byte) (string, error)
	BucketExist(bucketName string) bool
	Close() error
}
type SQLRepository interface {
	Query(query string, args ...any) (*sql.Rows, error)
	GetAllTables() ([]string, error)
	GetColumns(tableName string) ([]string, error)
	Close() error
}
