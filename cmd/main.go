package main

import (
	"context"
	"flag"
	"fmt"

	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	"log"
	"os"

	internal "github.com/Dmitry-dms/nirs/internal"
	"github.com/Dmitry-dms/nirs/internal/repository"
)

var(
	openBrowser bool
	port string
)

func init() {
	flag.BoolVar(&openBrowser, "b", false, "Открытие графического интерфейса при запуске")
	flag.StringVar(&port, "p", "8080", "Укажите порт, который прослушивает сервер. Если порт не стандартный, необходимо пересобрать графический интерфейс")
}

func main() {
	flag.Parse()

	logger := log.New(os.Stdout, "SYSTEM: ", 1)
	KVRepo := repository.NewRamCache()
	config := internal.Config{
		Address: fmt.Sprintf(":%s", port),
		KVRepo:  KVRepo,
		Logger:  logger,
	}

	c := internal.NewCore(config)
	// cat := c.GetCatalog(c.PerechenName)
	// fmt.Println(len(cat.Terrorists))
	// c.StoreAllKeys([]byte("bucketName []byte"), cat)

	closeCh := make(chan os.Signal, 1)
	signal.Notify(closeCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go c.StartServer()
	// Открытие вкладки в браузере по указанному адресу
	if openBrowser {
		open(fmt.Sprintf("http://localhost:%s/", port))
	}
	logger.Printf("Сайт доступен по адресу http://localhost%s/ \n", config.Address)
	logger.Println("Не закрывайте данное окно. Закрыть после окончания работы.")
	<-closeCh
	c.Shutdown(context.Background())
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
