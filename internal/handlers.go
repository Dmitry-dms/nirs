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
	mux.HandleFunc("/history", c.getHistory)
	mux.HandleFunc("/opt", c.getOptions)
	return mux
}
func (c *Core) searchHandler(rw http.ResponseWriter, r *http.Request) {
	var sr SearchRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&sr)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}

	sqlName := ""
	pName := ""
	//поиск названия sql базы и название перечня
	for _, v := range c.Settings.SQLOptions {
		if v.Path == sr.SQLPAth {
			sqlName = v.Name
		}
	}
	for _, v := range c.Settings.XMLOptions {
		if v.Path == sr.PPath {
			pName = v.Name
		}
	}

	c.getCatalog(sr.PPath)

	repo := repository.NewSqlite(sr.SQLPAth)
	res := c.Search(repo, sr.PPath, sr.Table, sr.Columns)
	repo.Close()

	history := HistoryElement{
		Date:    time.Now().Format("2006.01.02 15:04"),
		SQLName: sqlName,
		PName:   pName,
		ID:      sr.Uid,
		Columns: sr.Columns,
		Rows:    res,
	}
	c.History = append(c.History, history)
	data, err := json.Marshal(history)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Write(data)
}
// func (c *Core) showHandler(rw http.ResponseWriter, r *http.Request) {

// 	h, err := os.ReadFile("settings.json")
// 	if err != nil {
// 		return
// 	}
// 	s := string(h)
// 	fmt.Fprint(rw, s)
// }
// func (c *Core) write(rw http.ResponseWriter, r *http.Request) {
// 	var data *os.File
// 	var err error
// 	defer data.Close()
// 	data, err = os.OpenFile("settings.json", os.O_APPEND, 0666)
// 	if err != nil {
// 		data, err = os.Create("settings.json")
// 	}
// 	fmt.Fprint(data, "show page")
// }
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
	data, err := json.Marshal(c.Settings)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
	rw.Write(data)
}
