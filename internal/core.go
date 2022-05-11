package internal

import (
	"context"
	"fmt"

	"encoding/xml"

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
	KVRepo       repository.KVRepository
	SQLname      string
	Columns      []string `json:"columns"`
	colsEng      []string
	PerechenName string
	History      History
	logger       *log.Logger
	mu           *sync.Mutex
}

type Config struct {
	Address string
	KVRepo  repository.KVRepository
	Sqlite  *repository.SqliteRepository
	Logger  *log.Logger
}

func NewCore(c Config) *Core {
	core := Core{
		Server: http.Server{
			Addr: c.Address,
		},
		KVRepo: c.KVRepo,
		logger: c.Logger,
		mu:     &sync.Mutex{},
	}
	routes := initRoutes(&core)
	core.Server.Handler = routes

	cols, err := loadColumns("settings.json")
	if err != nil {
		core.logger.Fatal(err)
	}
	var colsRus []string
	var colsEng []string
	for _, v := range cols {
		colsRus = append(colsRus, v.Name)
		colsEng = append(colsEng, v.NameColumn)
	}
	core.Columns = colsRus
	core.colsEng = colsEng

	core.SQLname, core.PerechenName = core.getSettings()

	h, err := loadHistory("history.json")
	if err != nil {
		core.logger.Println("Файла history.json не найдено. Он будет создан после проведения первой проверки.")
	}

	core.GetCatalog(core.PerechenName)

	core.History = h
	return &core
}

func (c *Core) GetCatalog(path string) *Catalog {

	var catalogRaw XMLCatalog
	var t []*Terrorist

	c.ReadXMLFromDir(path, &catalogRaw)
	t = c.AggregateStructs(&catalogRaw, t)
	catalog := catalogRaw.ConvertCatalog(t)

	c.StoreAllKeys([]byte(path), &catalog)

	return &catalog
}

func (c *Core) StartServer() error {
	return c.ListenAndServe()
}

func (c *Core) Shutdown(ctx context.Context) {
	c.storeHistory()
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
			ter.Name = trimCDATASuffixAndPrefix(ter.Name)
			ter.BirthDate = trimCDATASuffixAndPrefix(ter.BirthDate)
			ter.Address = trimCDATASuffixAndPrefix(ter.Address)
			ter.Resolution = trimCDATASuffixAndPrefix(ter.Resolution)
			ter.BirthPlace = trimCDATASuffixAndPrefix(ter.BirthPlace)
			ter.Passport = trimCDATASuffixAndPrefix(ter.Passport)
			mu.Lock()
			t = append(t, ter.ConvertTerr(c.logger))
			mu.Unlock()
			wg.Done()
		}(ter, wg)
	}
	wg.Wait()

	return t
}

func (c *Core) SearchOne(value string) bool {
	return c.KVRepo.GetValue(value)
}
func (c *Core) Search(repo repository.SqliteRepository, tableName string) [][]Row[any] {
	rows, err := repo.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	pool := pool.NewPool(10, 3, 10)
	wg := &sync.WaitGroup{}
	var res [][]Row[any]
	mu := sync.Mutex{}
	for rows.Next() {
		wg.Add(1)
		var maps []map[string]any

		columns := make([]any, len(cols))
		columnPointers := make([]any, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		if err := rows.Scan(columnPointers...); err != nil {
			continue
		}

		m := make(map[string]any)
		for i, colName := range cols {
			val := columnPointers[i].(*any)
			m[colName] = *val
		}
		maps = append(maps, m)
		pool.Schedule(func() {
			r, isTrue := c.process(m, cols)
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
func (c *Core) process(m map[string]any, cols []string) ([]Row[any], bool) {
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
		ok := c.KVRepo.GetValue(v)
		// if err != nil {
		// 	c.logger.Println(err)
		// 	return nil, false
		// }
		if ok {
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

func (c *Core) GetValue(key string) bool {
	return c.KVRepo.GetValue(key)
}
func (c *Core) StoreAllKeys(bucketName []byte, catalog *Catalog) {
	// wg := &sync.WaitGroup{}
	for _, ter := range catalog.Terrorists {
		// wg.Add(1)
		// go func(ter *Terrorist, wg *sync.WaitGroup) {
		for _, n := range ter.Names {
			c.KVRepo.AddValue(n)
		}
		for _, a := range ter.Address {
			c.KVRepo.AddValue(a)
		}
		for _, p := range ter.Passport.SerialAndNum {
			c.KVRepo.AddValue(p)
		}
		c.KVRepo.AddValue(ter.INN)
		// wg.Done()
		// }(ter, wg)
	}
	// wg.Wait()
}
