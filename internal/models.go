package internal

import (
	"log"
	"regexp"
	"strings"
)

type History []HistoryElement

type HistoryElement struct {
	Date    string       `json:"date"`
	PName   string       `json:"pName"`
	ID      string       `json:"id"`
	Columns []string     `json:"columns"`
	Rows    [][]Row[any] `json:"rows"`
}

type SearchRequest struct {
	Uid string `json:"uid"`
}

type SearchOneRequest []SearchOneElement

type SearchOneElement struct {
	ID    int64  `json:"id"`
	Field string `json:"field"`
}
type SearchOneResponse struct {
	Id []int `json:"ids"`
}

type Row[T any] struct {
	Field    T    `json:"field"`
	Selected bool `json:"selected"`
}

type Columns []Column

type Column struct {
	Name       string `json:"name"`
	NameColumn string `json:"name_column"`
	Index      int64  `json:"index"`
}

type XMLCatalog struct {
	Num        string  `xml:"NUM,attr"`
	Date       string  `xml:"DATE,attr"`
	Id         string  `xml:"ID,attr"`
	Terrorists []*Terr `xml:"TERRORISTS"`
}
type Catalog struct {
	Num        string
	Date       string
	Id         string
	Terrorists []*Terrorist
}
type Terr struct {
	Name        string `xml:"TERRORISTS_NAME"`
	IsExtremist bool
	PersonType  int    `xml:"PERSON_TYPE"`
	INN         string `xml:"INN"`
	BirthDate   string `xml:"BIRTH_DATE"`
	Address     string `xml:"ADDRESS"`
	Resolution  string `xml:"TERRORISTS_RESOLUTION"`
	BirthPlace  string `xml:"BIRTH_PLACE"`
	Passport    string `xml:"PASSPORT"`
}

type Terrorist struct {
	Names       []string
	IsExtremist bool
	PersonType  int
	INN         string
	BirthDate   string
	Address     []string
	Resolution  string
	BirthPlace  string
	Passport    Passport
}
type Passport struct {
	RawData      []string
	SerialAndNum []string
}

func (c XMLCatalog) ConvertCatalog(terr []*Terrorist) Catalog {
	return Catalog{
		Num:        c.Num,
		Date:       c.Date,
		Id:         c.Id,
		Terrorists: terr,
	}
}
// func trimSpace(s []string) []string {
// 	var n []string
// 	for _, j := range s {
// 		m := strings.TrimSpace(j)
// 		n = append(n, m)
// 	}
// 	return n
// }
func convertInn(s string) string {
	inn := removeSpaces(s)
	if len(inn) == 0 {
		return "NULL"
	} else {
		return inn
	}
}
func (t *Terr) ConvertTerr(l *log.Logger) *Terrorist {
	names := splitNames(t.Name)
	adr := splitAddress(t.Address)
	pass := splitPassport(t.Passport, l)
	inn := convertInn(t.INN)
	return &Terrorist{
		Names:       names,
		IsExtremist: t.IsExtremist,
		PersonType:  t.PersonType,
		INN:         inn,
		BirthDate:   t.BirthDate,
		Address:     adr,
		Resolution:  t.Resolution,
		BirthPlace:  t.BirthPlace,
		Passport:    pass,
	}
}

func splitNames(s string) []string {
	var result []string
	pattern := "\\(+.*"
	r := regexp.MustCompile(pattern)
	secPart := r.FindAllString(s, 5) // ?????????? ???????????? ???? ?????? ?????????? ?? ?????????????? regex, ???????????????? ???????????? ?????????? ?? ????????????
	if secPart == nil {              // ???????? ?????????????? ???? ??????????????????, ???????????? ?????????? ?????????? ???????? ??????
		result = append(result, s)
		return result
	}
	secPart2 := secPart[0][1:(len(secPart[0]) - 1)] // ?????????????????????? ???? ???????????? ?????????????? ????????????
	splitted := strings.Split(secPart2, ";")
	for _, name := range splitted {
		k := strings.Split(name, "(") // ?????????? ???? ??????????????????
		if len(k) == 1 {
			result = append(result, removeSpaces(name))
		} else if len(k) == 3 {
			for _, j := range k {
				l := strings.Split(j, ")")
				for _, f := range l {
					if removeSpaces(f) == "" {
						continue
					} else {
						result = append(result, f)
					}
				}
			}
		} else {
			first := removeSpaces(k[0])
			second := removeSpaces(k[1][:len(k[1])-1])
			result = append(result, first, second)
		}
	}
	fistPart := s[:(len(s) - len(secPart2) - 2)]
	result = append(result, removeSpaces(fistPart))
	return removeDuplicateStr(result) // ?????????????????????? ???? ???????????????????? ??????????
}
func removeDuplicateStr(strSlice []string) []string {
	allKeys := make(map[string]bool)
	list := []string{}
	for _, item := range strSlice {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
func removeSpaces(s string) string {
	return strings.TrimSpace(s)
}
func splitAddress(s string) []string {
	var a []string
	addresses := strings.Split(s, ";")
	for _, j := range addresses {
		gg := strings.TrimSuffix(j, ",")
		a = append(a, removeSpaces(gg))
	}
	return a
}
func splitPassport(s string, l *log.Logger) Passport {
	if len(s) <= 10 {
		return Passport{
			RawData:      nil,
			SerialAndNum: nil,
		}
	} else {
		var singlePass []string
		var serial []string
		raw := strings.Split(s[:len(s)-1], ",")
		singlePass = append(singlePass, raw...)
		for _, p := range singlePass {
			if len(p) < 9 { // ???????? ???????????????????? ?????????? ?????????? ???? ???? ??????????????, ???????????????????? ????????????????????
				l.Printf("?????????????? ???????????????? ????????????????, ???????????????????????? ????????????: [%#v] [%#v]", raw[0], raw[1])
				return Passport{
					RawData:      raw,
					SerialAndNum: nil,
				}
			}
			runeString := []rune(p)
			if runeString[8] == '??' {
				serial = append(serial, string(runeString[12:23]))
			} else if runeString[8] == '??' {
				serial = append(serial, string(runeString[14:27]))
			} else if runeString[10] == 'A' {
				serial = append(serial, string(runeString[10:22]))
			} else {
				serial = append(serial, strings.TrimSpace(string(runeString[10:20])))
			}
		}
		return Passport{
			RawData:      singlePass,
			SerialAndNum: serial,
		}
	}
}
func trimCDATASuffixAndPrefix(s string) string {
	var l string
	if len(s) <= 13 {
		l = "NULL"
	} else {
		n := s[10:]
		n = strings.ReplaceAll(n, "*", "")
		l = n[:len(n)-4]
	}
	return l
}
