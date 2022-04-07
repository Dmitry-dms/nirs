package internal

import (
	"context"

	"encoding/xml"
	"fmt"

	"log"
	"net/http"

	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/Dmitry-dms/nirs/internal/repository"
	"github.com/Dmitry-dms/nirs/pkg/pool"
)

type Core struct {
	http.Server
	KVRepo   repository.KVRepository
	Sqlite   *repository.SqliteRepository
	Settings Options
	History  History
	logger   *log.Logger
	mu       *sync.Mutex
}

type Config struct {
	Address string
	KVRepo  repository.KVRepository
	Sqlite  *repository.SqliteRepository
	Logger  *log.Logger
}

var aggregatedCatalogs = make(map[string]*Catalog, 10)

func NewCore(c Config) *Core {
	core := Core{
		Server: http.Server{
			Addr: c.Address,
		},
		KVRepo: c.KVRepo,
		logger: c.Logger,
		mu:     &sync.Mutex{},
		Sqlite: c.Sqlite,
	}
	routes := initRoutes(&core)
	core.Server.Handler = routes
	s := core.getSettings()
	// s, _ := loadSettings("settings.json")
	h, err := loadHistory("history.json")
	if err != nil {
		fmt.Println("failed to load history")
	}
	core.Settings = s
	core.History = h
	return &core
}

func (c *Core) getCatalog(path string) *Catalog {
	//проверка на существующий каталог
	if v, ok := aggregatedCatalogs[path]; ok == true {
		return v
	}
	var catalogRaw XMLCatalog
	var t []*Terrorist
	c.ReadXMLFromDir(path, &catalogRaw)
	c.AggregateStructs(&catalogRaw, t)
	catalog := catalogRaw.ConvertCatalog(t)
	c.StoreAllKeys([]byte(path), &catalog)
	aggregatedCatalogs[path] = &catalog
	return &catalog
}

func (c *Core) StartServer() error {
	return c.ListenAndServe()
}

func (c *Core) Shutdown(ctx context.Context) {
	// c.storeSettings()
	c.storeHistory()
	c.KVRepo.Close()
	c.Sqlite.Close()
	c.Server.Shutdown(ctx)
}
func (c *Core) ReadXMLFromDir(path string, cat *XMLCatalog) error {
	xmlFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer xmlFile.Close()
	byteValue, err := ioutil.ReadAll(xmlFile)
	xml.Unmarshal(byteValue, &cat)
	return err
}

func (c *Core) AggregateStructs(cat *XMLCatalog, t []*Terrorist) []*Terrorist {
	wg := &sync.WaitGroup{}

	mu := &sync.Mutex{}
	for _, ter := range cat.Terrorists {
		wg.Add(1)
		go func(ter *Terr, wg *sync.WaitGroup) {
			if strings.Contains(ter.Name, "*") {
				ter.IsExtremist = true
			}
			ter.Name = TrimSuffixAndPrefix(ter.Name)
			ter.BirthDate = TrimSuffixAndPrefix(ter.BirthDate)
			ter.Address = TrimSuffixAndPrefix(ter.Address)
			ter.Resolution = TrimSuffixAndPrefix(ter.Resolution)
			ter.BirthPlace = TrimSuffixAndPrefix(ter.BirthPlace)
			ter.Passport = TrimSuffixAndPrefix(ter.Passport)
			mu.Lock()
			t = append(t, ter.ConvertTerr(c.logger))
			mu.Unlock()
			wg.Done()
		}(ter, wg)
	}
	wg.Wait()

	return t
}


func (c *Core) Search(bucketName, tableName string, cols []string) [][]Row[any] {
	rows, err := c.Sqlite.Db.Query(fmt.Sprintf("SELECT %s FROM %s", strings.Join(cols, ", "), tableName))
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	pool := pool.NewPool(10, 3, 10)
	wg := &sync.WaitGroup{}
	var res [][]Row[any]
	mu := sync.Mutex{}
	for rows.Next() {
		wg.Add(1)
		var maps []map[string]any

		columns := make([]any, len(cols))
		columnPointers := make([]any, len(cols))
		for i, _ := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			continue
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]any)
		for i, colName := range cols {
			val := columnPointers[i].(*any)
			m[colName] = *val
		}
		maps = append(maps, m)
		pool.Schedule(func() {
			r, isTrue := c.process(bucketName, m, cols)
			if isTrue != true {
				wg.Done()
				return
			} else {
				mu.Lock()
				res = append(res, r)
				mu.Unlock()
				wg.Done()
			}
		})
	}
	wg.Wait()
	return res
}
func (c *Core) process(bucketName string, m map[string]any, cols []string) ([]Row[any],bool) {
	numCols := len(cols)
	numAfter := numCols

	row := make([]Row[any], numCols)
	colIndexes := make(map[string]int, numCols)
	//находим индекс колонки в массиве
	for i, v := range cols {
		colIndexes[v] = i
	}
	for k, v := range m {
		v := fmt.Sprintf("%v", v)
		var r Row[any]
		one, _ := c.KVRepo.GetValue([]byte(bucketName), []byte(v))
		if one == "true" {
			r.Selected = true
			numAfter--
		} else {
			r.Selected = false
			
		}
		r.Field = v
		colIndex := colIndexes[k]
		row[colIndex] = r
	}
	if numCols == numAfter {
		return nil, false
	}
	return row, true

}

func (c *Core) GetValue(bucketName, key []byte) (string, error) {
	return c.KVRepo.GetValue(bucketName, key)
}
func (c *Core) StoreAllKeys(bucketName []byte, catalog *Catalog) {
	wg := &sync.WaitGroup{}
	for _, ter := range catalog.Terrorists {
		wg.Add(1)
		go func(ter *Terrorist, wg *sync.WaitGroup) {
			for _, n := range ter.Names {
				c.KVRepo.AddValue(bucketName, []byte(n), []byte("true"))
			}
			for _, a := range ter.Address {
				c.KVRepo.AddValue(bucketName, []byte(a), []byte("true"))
			}
			for _, p := range ter.Passport.SerialAndNum {
				c.KVRepo.AddValue(bucketName, []byte(p), []byte("true"))
			}
			c.KVRepo.AddValue(bucketName, []byte(ter.INN), []byte("true"))
			wg.Done()
		}(ter, wg)
	}
	wg.Wait()
}
