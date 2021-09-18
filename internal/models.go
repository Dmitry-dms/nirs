package internal

import (
	"strings"
)

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
type PassType int

const (
	Russia PassType = iota
	Ussr
	Other
)

type Passport struct {
	RawData      []string
	SerialAndNum []string
	Type         PassType
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

func (c XMLCatalog) ConvertCatalog(terr []*Terrorist) Catalog {
	return Catalog{
		Num:        c.Num,
		Date:       c.Date,
		Id:         c.Id,
		Terrorists: terr,
	}
}
func trimSpace(s []string) []string {
	var n []string
	for _, j := range s {
		m := strings.TrimSpace(j)
		n = append(n, m)
	}
	return n
}
func (t *Terr) ConvertTerr() *Terrorist {
	names := splitNames(t.Name)
	adr := splitAddress(t.Address)
	pass := splitPassport(t.Passport)
	return &Terrorist{
		Names:       names,
		IsExtremist: t.IsExtremist,
		PersonType:  t.PersonType,
		INN:         t.INN,
		BirthDate:   t.BirthDate,
		Address:     adr,
		Resolution:  t.Resolution,
		BirthPlace:  t.BirthPlace,
		Passport:    pass,
	}
}
func splitNames(s string) []string {
	var names []string
	secondName := strings.Split(s, "(")
	if len(secondName) == 1 {
		names = append(names, s)
	} else {
		firstName := s[:(len(s) - len(secondName[1]) - 1)]
		n := strings.ReplaceAll(secondName[1], ")", "")
		newName := strings.Split(n, ";")
		names = append(names, firstName)
		newName = trimSpace(newName)
		names = append(names, newName...)
	}
	return names
}
func splitAddress(s string) []string {
	var a []string
	addresses := strings.Split(s, ";")
	a = append(a, addresses...)
	return a
}
func splitPassport(s string) Passport {
	if len(s) <= 35 {
		return Passport{
			RawData:      nil,
			SerialAndNum: nil,
			Type: 10,
		}
	} else {
		var r []string
		var serial []string
		raw := strings.Split(s[:len(s)-1], ",")
		r = append(r, raw...)
		for _, p := range r { //13-25
			serial = append(serial, p[21:32])
		}
		return Passport{
			RawData:      r,
			SerialAndNum: serial,
		}
	}
}
func TrimSuffixAndPrefix(s string) string {
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