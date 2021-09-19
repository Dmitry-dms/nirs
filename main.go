package main

import (
	"fmt"
	"strings"
)

type n struct {
	T []string
}

func main() {
	test := "ПАСПОРТ СССР: III-ОЖ 549177 ВЫДАН ЛЕНИНСКИМ ОВД Г. ОРДЖОНИКИДЗЕ 23.05.1977, PASSPORT: 11 ВА18320 ВЫДАН ПАСПОРТ ГРАЖДАНИНА ГРУЗИИ 21.04.2014,"
	final := SplitPassport(test)
	fmt.Println(final.SerialAndNum)
	//fmt.Println(string([]rune(test)[13:25]))
}
func splitAddress(s string) []string {
	var a []string
	addresses := strings.Split(s, ";")
	a = append(a, addresses...)
	return a
}



type Passport struct {
	RawData      []string
	SerialAndNum []string
}

func SplitPassport(s string) Passport {
	if s == "" {
		return Passport{
			RawData:      nil,
			SerialAndNum: nil,
		}
	} else {
		var r []string
		var serial []string
		raw := strings.Split(s[:len(s)-1], ",")
		r = append(r, raw...)
		for _, p := range r { 
			r := []rune(p)
			if r[8] == 'Р' {
				serial = append(serial, string(r[12:23]))
			} else if r[8] == 'С' {
				serial = append(serial, string(r[14:27]))
			} else {
				serial = append(serial, strings.TrimSpace(string(r[10:20])))
			}
		}
		return Passport{
			RawData:      r,
			SerialAndNum: serial,
		}
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
