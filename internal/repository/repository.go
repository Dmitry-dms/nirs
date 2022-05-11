package repository

import "database/sql"

type KVRepository interface {
	AddValue(key string)
	GetValue(key string) bool
}
type SQLRepository interface {
	Query(query string, args ...any) (*sql.Rows, error)
	GetAllTables() ([]string, error)
	GetColumns(tableName string) ([]string, error)
	Close() error
}
