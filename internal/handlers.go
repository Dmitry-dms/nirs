package internal

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Dmitry-dms/nirs/internal/repository"
)

func initRoutes(c *Core) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("dist")))
	mux.HandleFunc("/search", c.searchHandler)
	mux.HandleFunc("/search-one", c.searchOneHandler)
	mux.HandleFunc("/history", c.getHistory)
	mux.HandleFunc("/opt", c.getOptions)
	return mux
}
func (c *Core) searchOneHandler(rw http.ResponseWriter, r *http.Request) {
	var sor SearchOneRequest
	
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&sor)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	resId := make([]int, 0)
	for _, v := range sor {
		if c.SearchOne(v.Field) {
			resId = append(resId, int(v.ID))
		}
	}
	response := SearchOneResponse{
		Id: resId,
	}
	data, err := json.Marshal(response)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Write(data)
}
func (c *Core) searchHandler(rw http.ResponseWriter, r *http.Request) {
	var sr SearchRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&sr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	c.GetCatalog(c.PerechenName)

	repo := repository.NewSqlite(c.SQLname)
	res := c.Search(repo, "persons")
	repo.Close()


	history := HistoryElement{
		Date:    time.Now().Format("2006.01.02 15:04"),
		PName:   c.PerechenName,
		ID:      sr.Uid,
		Columns: c.Columns,
		Rows:    res,
	}
	c.History = append(c.History, history)

	data, err := json.Marshal(history)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Write(data)
}

func (c *Core) getHistory(rw http.ResponseWriter, r *http.Request) {
	var file *os.File
	var err error
	defer file.Close()
	file, err = os.OpenFile("history.json", os.O_APPEND, 0666)
	if err != nil {
		file, err = os.Create("history.json")
	}
	data, err := json.Marshal(c.History)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Write(data)
}

func (c *Core) getOptions(rw http.ResponseWriter, r *http.Request) {
	data, err := json.Marshal(c.Columns)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Write(data)
}
