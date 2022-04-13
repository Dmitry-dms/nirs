package main

import (
	"context"


	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	"log"
	"os"

	internal "github.com/Dmitry-dms/nirs/internal"
	"github.com/Dmitry-dms/nirs/internal/repository"
)



func main() {

	logger := log.New(os.Stdout, "SYSTEM: ", 1)
	KVRepo := repository.NewBoltDB("perechen.db")
	config := internal.Config{
		Address: ":8080",
		KVRepo:  KVRepo,
		Logger:  logger,
	}
	core := internal.NewCore(config)

	closeCh := make(chan os.Signal, 1)
	signal.Notify(closeCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go core.StartServer()
	go open("http://localhost:8080/")
	logger.Printf("Сайт доступен по адресу http://localhost%s/ \n", config.Address)
	logger.Println("Не закрывайте данное окно. Закрыть после окончания работы.")
	<-closeCh
	core.Shutdown(context.Background())
}

func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}
