package main

import (
	//"fmt"
	"log"
	"os"
	"sync"

	internal "github.com/Dmitry-dms/nirs/internal"
	"github.com/Dmitry-dms/nirs/internal/repository"
)

func main() {
	var catalogRaw internal.XMLCatalog
	var t []*internal.Terrorist
	logger := log.New(os.Stdout, "SYSTEM: ", 1)
	KVRepo := repository.NewBoltDB("test.db")
	core := internal.NewCore(KVRepo, logger)
	err := core.ReadXMLFromDir("ter.xml", &catalogRaw)
	if err != nil {
		log.Println(err)
	}
	
	wg := &sync.WaitGroup{}
	core.AggregateStructs(&catalogRaw, t, wg)	
	wg.Wait()

	//catalog := catalogRaw.ConvertCatalog(t)
	// for i := 0; i < 20; i++ {
	// 	fmt.Println(i)
	// 	fmt.Println(newCatalog.Terrorists[i].Names)
	// 	fmt.Println("-------------")
	// }
	// core.StoreAllKeys(&newCatalog)
	//fmt.Println(catalogRaw.Terrorists[0].Name)
}
