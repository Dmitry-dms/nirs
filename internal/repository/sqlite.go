package repository

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteRepository struct {
	Db *sql.DB
}

func NewSqlite(name string) SqliteRepository {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		log.Fatalf("Невозможно открыть базу данных %s: %v", name, err)
	}
	return SqliteRepository{Db: db}
}
func (r *SqliteRepository) Close() error {
	return r.Db.Close()
}

func (r *SqliteRepository) Query(query string, args ...any) (*sql.Rows, error) {
	return r.Db.Query(query, args...)
}

func (r *SqliteRepository) GetAllTables() ([]string, error) {
	query, err := r.Db.Query("SELECT name FROM sqlite_schema WHERE type ='table' AND name NOT LIKE 'sqlite_%'")
	if err != nil {
		return nil, err
	}
	var tables []string
	for query.Next() {
		var table string
		err := query.Scan(&table)
		if err != nil {
			fmt.Println(err)
			continue
		}
		tables = append(tables, table)
	}
	return tables, nil
}
func (r *SqliteRepository) GetColumns(tableName string) ([]string, error) {
	query, err := r.Db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return nil, err
	}
	return query.Columns()
}
