package internal

import (
	"encoding/json"
	"path/filepath"
	"strings"
	"os"
)

func loadColumns(path string) (Columns, error) {
	var s Columns
	file, err := os.ReadFile(path)
	if err != nil {
		return Columns{}, err
	}
	err = json.Unmarshal(file, &s)
	if err != nil {
		return Columns{}, err
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

func (c *Core) getSettings() (sqlName, perechenName string) {
	sqlName = searchFiles(".db")
	perechenName = searchFiles(".xml")
	return
}

func searchFiles(extFormat string) string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exPath := filepath.Dir(ex)
	files, _ := os.ReadDir(exPath)
	for _, name := range files {
		temp := name.Name()
		if !strings.Contains(temp, extFormat) {
			continue
		}
		runes := []rune(temp)
		s := runes[len(runes)-len([]rune(extFormat)):]
		if string(s) == extFormat {
			return temp
		}
	}
	return ""
}
