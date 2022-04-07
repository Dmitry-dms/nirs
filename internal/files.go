package internal

import (
	"encoding/json"
	"io"
	"path/filepath"
	"strings"

	"log"
	"os"
)

func loadSettings(path string) (Options, error) {
	var s Options
	file, err := os.ReadFile(path)
	if err != nil {
		log.Printf("there is no such file: %v", err)
		return Options{}, err
	}
	err = json.Unmarshal(file, &s)
	if err != nil {
		log.Printf("unmarshaling failed: %v", err)
		return Options{}, err
	}
	return s, nil
}
func loadHistory(path string) (History, error) {
	var h History
	file, err := os.ReadFile(path)
	if err != nil {
		log.Printf("there is no such file: %v", err)
		return nil, err
	}
	err = json.Unmarshal(file, &h)
	if err != nil {
		log.Printf("unmarshaling failed: %v", err)
		return nil, err
	}
	return h, nil
}

func (c *Core) storeSettings() error {
	temp, err := os.CreateTemp("/", "temp.file")
	if err != nil {
		log.Printf("failed to create temp file: %v", err)
		return err
	}
	var data *os.File
	defer data.Close()
	data, err = os.OpenFile("settings.json", os.O_APPEND, 0666)
	if err != nil {
		data, err = os.Create("settings.json")
		log.Printf("error while write: %v", err)
	}
	io.Copy(temp, data)
	json, err := json.Marshal(c.Settings)
	if err != nil {
		log.Printf("marshaling failed: %v", err)
		return err
	}
	_, err = data.Write(json)
	return err
}
func (c *Core) storeHistory() error {
	temp, err := os.CreateTemp("/", "htemp.file")
	if err != nil {
		log.Printf("failed to create temp file: %v", err)
		return err
	}
	var data *os.File
	defer data.Close()
	data, err = os.OpenFile("history.json", os.O_APPEND, 0666)
	if err != nil {
		data, err = os.Create("history.json")
		log.Printf("error while write: %v", err)
	}
	io.Copy(temp, data)
	json, err := json.Marshal(c.History)
	if err != nil {
		log.Printf("marshaling failed: %v", err)
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
		tables, err := c.Sqlite.GetAllTables()
		if err != nil {
			c.logger.Println(err)
		}
		tab := make([]Table, len(tables))
		for i, table := range tables {
			cols, _ := c.Sqlite.GetColumns(table)
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
