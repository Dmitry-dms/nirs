package main

import (
	//"fmt"
	"fmt"
	"log"
	"os"
	"time"

	internal "github.com/Dmitry-dms/nirs/internal"
	"github.com/Dmitry-dms/nirs/internal/repository"
)

func main() {
	now := time.Now()
	var catalogRaw internal.XMLCatalog
	var t []*internal.Terrorist
	logger := log.New(os.Stdout, "SYSTEM: ", 1)
	KVRepo := repository.NewBoltDB("test.db")
	core := internal.NewCore(KVRepo, logger)
	err := core.ReadXMLFromDir("ter.xml", &catalogRaw)
	if err != nil {
		log.Println(err)
	}

	t = core.AggregateStructs(&catalogRaw, t)

	catalog := catalogRaw.ConvertCatalog(t)

	// for i := 0; i < 2; i++ {
	// 	fmt.Println(i)
	// 	fmt.Println(catalog.Terrorists[i].Names)
	// 	fmt.Println("-------------")
	// }
	core.StoreAllKeys(&catalog)
	fmt.Printf("Длина изначально - %d \n", len(catalogRaw.Terrorists))
	fmt.Printf("Длина после конвертации - %d \n", len(t))
	fmt.Printf("Время выполнения - %s \n", time.Since(now))
}

