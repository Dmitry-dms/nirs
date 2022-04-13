package internal

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"os"

	"github.com/Dmitry-dms/nirs/internal/repository"
)

func loadSettings(path string) (Options, error) {
	var s Options
	file, err := os.ReadFile(path)
	if err != nil {
		return Options{}, err
	}
	err = json.Unmarshal(file, &s)
	if err != nil {
		return Options{}, err
	}
	return s, nil
}
func loadHistory(path string) (History, error) {
	var h History
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, &h)
	if err != nil {
		return nil, err
	}
	return h, nil
}


func (c *Core) storeHistory() error {	
	data, err := os.Create("history.json")
	if err != nil {
		return err
	}
	defer data.Close()
	json, err := json.Marshal(c.History)
	if err != nil {
		return err
	}
	_, err = data.Write(json)
	return err
}

func (c *Core) getSettings() Options {
	sql := searchFiles(".db")
	xml := searchFiles(".xml")
	sqls := make([]SQLOption, len(sql))
	xmls := make([]XMLOption, len(xml))
	for i, v := range sql {
		name := []rune(v)[:len([]rune(v))-3]
		sqLite := repository.NewSqlite(v)
		tables, err := sqLite.GetAllTables()
		if err != nil {
			c.logger.Println(err)
		}
		tab := make([]Table, len(tables))
		for i, table := range tables {
			cols, _ := sqLite.GetColumns(table)
			t := Table{
				Name:    table,
				Columns: cols,
			}
			tab[i] = t
		}

		sqlOpt := SQLOption{
			Name: string(name),
			Path: v,
			Tables: tab,
		}
		sqls[i] = sqlOpt
		sqLite.Close()
	}
	for i, v := range xml {
		name := []rune(v)[:len([]rune(v))-4]
		xmlOpt := XMLOption{
			Name: string(name),
			Path: v,
		}
		xmls[i] = xmlOpt
	}

	opt := Options{
		SQLOptions: sqls,
		XMLOptions: xmls,
	}
	return opt

}

func searchFiles(extFormat string) []string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	var result []string
	exPath := filepath.Dir(ex)
	files, _ := os.ReadDir(exPath)
	for _, name := range files {
		temp := name.Name()
		if !strings.Contains(temp, extFormat) {
			continue
		}
		//пропустим базу с перечнем
		if temp == "perechen.db" {
			continue
		}
		runes := []rune(temp)
		s := runes[len(runes)-len([]rune(extFormat)):]
		if string(s) == extFormat {
			result = append(result, temp)
		}
	}
	return result
}
