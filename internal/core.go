package internal

import (
	"encoding/xml"
	"fmt"
	"log"
	"time"

	"hash"
	"hash/fnv"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/Dmitry-dms/nirs/internal/repository"
	"github.com/Dmitry-dms/nirs/pkg/pool"
)

type Core struct {
	KVRepo repository.KVRepository
	//Sqlite repository.SQLRepository
	Sqlite *repository.SqliteRepository
	hash   hash.Hash32
	logger *log.Logger
	mu     *sync.Mutex
}

func NewCore(r repository.KVRepository, sql *repository.SqliteRepository, logger *log.Logger) *Core {
	g := fnv.New32()
	return &Core{
		KVRepo: r,
		logger: logger,
		hash:   g,
		mu:     &sync.Mutex{},
		Sqlite: sql,
	}
}
func (c *Core) Shutdown() {
	c.KVRepo.Close()
	c.Sqlite.Close()
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
	now := time.Now()
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
	fmt.Println(time.Since(now))
	return t
}

type Result struct {
	repository.SqlitePerson
	Res []bool
}

func (c *Core) Search(tableName string) []*Result {
	rows, err := c.Sqlite.Db.Query(fmt.Sprintf("select * from %s", tableName))
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	pool := pool.NewPool(10, 3, 10)
	wg := &sync.WaitGroup{}
	var res []*Result
	mu := sync.Mutex{}
	for rows.Next() {
		wg.Add(1)
		p := repository.SqlitePerson{}
		err := rows.Scan(&p.Id, &p.Name, &p.Passport, &p.Inn, &p.Address)
		if err != nil {
			log.Println(err)
		}
		pool.Schedule(func() {
			r := c.process(&p)
			if r == nil {
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
func (c *Core) process(r *repository.SqlitePerson) *Result {
	var h []bool
	one, _ := c.KVRepo.GetValue([]byte(r.Name))
	if one == "true" {
		h = append(h, true)
	} else {
		h = append(h, false)
	}
	two, _ := c.KVRepo.GetValue([]byte(r.Passport))
	if two == "true" {
		h = append(h, true)
	} else {
		h = append(h, false)
	}
	three, _ := c.KVRepo.GetValue([]byte(r.Inn))
	if three == "true" {
		h = append(h, true)
	} else {
		h = append(h, false)
	}
	four, _ := c.KVRepo.GetValue([]byte(r.Address))
	if four == "true" {
		h = append(h, true)
	} else {
		h = append(h, false)
	}

	if isTrue(h) {
		return &Result{
			Res:          h,
			SqlitePerson: *r,
		}
	} else {
		return nil
	}
}
func isTrue(s []bool) bool {
	for _, h := range s {
		if h == true {
			return true
		}
	}
	return false
}
func (c *Core) GetValue(key []byte) (string, error) {
	return c.KVRepo.GetValue(key)
}
func (c *Core) StoreAllKeys(catalog *Catalog) {
	wg := &sync.WaitGroup{}
	for _, ter := range catalog.Terrorists {
		wg.Add(1)
		go func(ter *Terrorist, wg *sync.WaitGroup) {
			for _, n := range ter.Names {
				c.KVRepo.AddValue([]byte(n), []byte("true"))
			}
			for _, a := range ter.Address {
				c.KVRepo.AddValue([]byte(a), []byte("true"))
			}
			for _, p := range ter.Passport.SerialAndNum {
				c.KVRepo.AddValue([]byte(p), []byte("true"))
			}
			c.KVRepo.AddValue([]byte(ter.INN), []byte("true"))
			wg.Done()
		}(ter, wg)
	}
	wg.Wait()
}
