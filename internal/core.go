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
//	"github.com/Dmitry-dms/nirs/pkg/pool"
)

type Core struct {
	KVRepo repository.KVRepository
	Sqlite repository.SQLRepository
	hash   hash.Hash32
	logger *log.Logger
	mu     *sync.Mutex
}

func NewCore(r repository.KVRepository, sql repository.SQLRepository, logger *log.Logger) *Core {
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

func (c *Core) Search(out chan Result, quit chan int) []*Result{
	rows := c.Sqlite.GetAllValues()
	mu := &sync.Mutex{}
	var s []*Result
	//pool := pool.NewPool(1000, 100, 1000)
	wg := &sync.WaitGroup{}
	for _, row := range rows {
		wg.Add(1)
		//pool.Schedule(
			func() {
				r := c.process(wg, row, out)
				if r == nil {
					return
				} else {
				mu.Lock()
				s=append(s, r)
				mu.Unlock()
				}
			}()
	//	)

	}
	wg.Wait()
	return s
	//quit <- 0
}
func (c *Core) process(w *sync.WaitGroup, r *repository.SqlitePerson, out chan Result) *Result {
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
	//res := [4]string{one, two, three, four}
	if isTrue(h) {
		s := Result{
			Res:          h,
			SqlitePerson: *r,
		}
		//out <- s
		w.Done()
		return &s
	} else {
		w.Done()
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
