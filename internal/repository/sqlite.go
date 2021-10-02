package repository

import (
	"database/sql"
	"log"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type SqliteRepository struct {
	Db *sql.DB
}
type SqlitePerson struct {
	Id int
	Name string
	Passport string
	Inn string
	Address string
}
func NewSqlite(name string) *SqliteRepository {
	db, err := sql.Open("sqlite3", name)
	if err != nil {
		log.Fatalf("Failed to open connection with sqlite: %v", err)
	}
	return &SqliteRepository{Db: db}
}
func (r *SqliteRepository) Close() error {
	return r.Db.Close()
}
func (r *SqliteRepository) GetAllValues() []*SqlitePerson {
	rows, err := r.Db.Query("select * from MOCKDATA")
    if err != nil {
        panic(err)
    }
    defer rows.Close()
	var persons []*SqlitePerson
	//w := &sync.WaitGroup{}
	mu := sync.Mutex{}
	for rows.Next() {
		//w.Add(1)
		//go func(wg *sync.WaitGroup) {
			p := SqlitePerson{}
			err := rows.Scan(&p.Id, &p.Name, &p.Passport, &p.Inn, &p.Address)
			if err != nil{
				log.Println(err)
			}
			mu.Lock()
			persons = append(persons, &p)
			mu.Unlock()
			//w.Done()
		//}(w)
	}
	//w.Wait()
	return persons
}