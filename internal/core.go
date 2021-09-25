package internal

import (
	"encoding/binary"
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
)

type Core struct {
	KVRepo repository.KVRepository
	hash   hash.Hash32
	logger *log.Logger
	mu     *sync.Mutex
}

func NewCore(r repository.KVRepository, logger *log.Logger) *Core {
	g := fnv.New32()
	return &Core{
		KVRepo: r,
		logger: logger,
		hash:   g,
		mu:     &sync.Mutex{},
	}
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
func (c *Core) GenerateHash(s string) []byte {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.hash.Write([]byte(s))
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, c.hash.Sum32())
	return bs
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
			wg.Done()
		}(ter, wg)
	}

	wg.Wait()
}
