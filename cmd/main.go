package main

import (
	// "context"
	"context"

	"log"
	"os"


	internal "github.com/Dmitry-dms/nirs/internal"
	"github.com/Dmitry-dms/nirs/internal/repository"
)

func main() {
	// now := time.Now()
	// var catalogRaw internal.XMLCatalog
	// var t []*internal.Terrorist

	logger := log.New(os.Stdout, "SYSTEM: ", 1)
	KVRepo := repository.NewBoltDB("perechen.db")
	Sqlite := repository.NewSqlite("people.db")
	config := internal.Config{
		Address: ":8080",
		KVRepo:  KVRepo,
		Sqlite:  Sqlite,
		Logger:  logger,
	}
	core := internal.NewCore(config)

	// err := core.ReadXMLFromDir("ter.xml", &catalogRaw)
	// if err != nil {
	// 	log.Println(err)
	// }

	// t = core.AggregateStructs(&catalogRaw, t)

	// catalog := catalogRaw.ConvertCatalog(t)
	// logger.Printf("Парсинг занял %s, количество записей - %d", time.Since(now), len(catalog.Terrorists))
	// closeCh := make(chan os.Signal, 1)
	// signal.Notify(closeCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	// <-closeCh

	// core.StoreAllKeys([]byte("names"), &catalog)
	//logger.Println("server starting")
	core.StartServer()
	defer core.Shutdown(context.Background())
	// cols := []string{"id","fio","passport","inn_ogrn","address"}
	// s := core.Search("names","MOCKDATA",cols)

	// // //fmt.Printf("Dlina - %d \n", len(s))
	// for _, j := range s {
	// 	// fmt.Printf("[%s;%s;%s;%s]\n", j.Name, j.Passport, j.Inn, j.Address)
	// 	fmt.Printf("%v \n", j)
	// }
	// // //core.Shutdown(context.Background())
	// fmt.Printf("Кол-во совпадений - %d \n", len(s))
	// fmt.Printf("Время выполнения - %s \n", time.Since(now))
}
