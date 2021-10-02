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
	//	var catalogRaw internal.XMLCatalog
	//var t []*internal.Terrorist
	logger := log.New(os.Stdout, "SYSTEM: ", 1)
	KVRepo := repository.NewBoltDB("test.db")
	Sqlite := repository.NewSqlite("people.db")
	core := internal.NewCore(KVRepo, Sqlite, logger)
	// err := core.ReadXMLFromDir("ter.xml", &catalogRaw)
	// if err != nil {
	// 	log.Println(err)
	// }

	//t = core.AggregateStructs(&catalogRaw, t)

	//catalog := catalogRaw.ConvertCatalog(t)
	// closeCh := make(chan os.Signal, 1)
	// signal.Notify(closeCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// <-closeCh

	// for i := 0; i < 2; i++ {
	// 	fmt.Println(i)
	// 	fmt.Println(catalog.Terrorists[i].Names)
	// 	fmt.Println("-------------")
	// }
	//core.StoreAllKeys(&catalog)
	// cn := make(chan internal.Result)
	// quit := make(chan int)
	//var res []internal.Result
	s := core.Search()

	// ff:
	// 	for {
	// 		select {
	// 		case r := <-cn:
	// 			res = append(res, r)
	// 		case <-quit:
	// 			break ff
	// 		}
	// 	}
	//fmt.Printf("Dlina - %d \n", len(s))
	// for _, j := range s {
	// 	fmt.Printf("[%s;%s;%s;%s]\n", j.Name, j.Passport, j.Inn, j.Address)
	// 	fmt.Printf("%v \n", j.Res)
	// }
	core.Shutdown()
	fmt.Printf("Кол-во совпадений - %d \n", len(s))
	//fmt.Printf("Длина изначально - %d \n", len(catalogRaw.Terrorists))
	//fmt.Printf("Длина после конвертации - %d \n", len(catalog.Terrorists))
	fmt.Printf("Время выполнения - %s \n", time.Since(now))
}
