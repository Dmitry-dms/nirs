package internal

import (
	"encoding/xml"
	"fmt"
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
}

func NewCore(r repository.KVRepository) *Core {
	g := fnv.New32()
	return &Core{KVRepo: r,
		hash: g}
}

func (c *Core) ReadXMLFromDir(path string, cat *XMLCatalog) error {
	xmlFile, err := os.Open("ter.xml")
	if err != nil {
		return err
	}
	defer xmlFile.Close()
	byteValue, err := ioutil.ReadAll(xmlFile)
	var catalog XMLCatalog
	xml.Unmarshal(byteValue, &catalog)
	return err
}

func (c *Core) AggregateStructs(cat *XMLCatalog, t []*Terrorist, wg *sync.WaitGroup) {
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
			t = append(t, ter.ConvertTerr())
			fmt.Println(ter)
			wg.Done()
		}(ter, wg)
	}
}
func (c *Core) generateHash(s string) (int, error) {
	return c.hash.Write([]byte(s))
}
func (c *Core) StoreAllKeys(catalog *Catalog) {
	wg := &sync.WaitGroup{}
	for _, ter := range catalog.Terrorists {
		wg.Add(1)
		go func(ter *Terrorist, wg *sync.WaitGroup) {
			for _, n := range ter.Names {
				c.KVRepo.AddValue(n, "true")
			}
			for _, a := range ter.Address {
				c.KVRepo.AddValue(a, "true")
			}
			for _, p := range ter.Passport.SerialAndNum {
				c.KVRepo.AddValue(p, "true")
			}
			wg.Done()
		}(ter, wg)
	}
	wg.Wait()
}
