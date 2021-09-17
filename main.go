package main

import (
	"fmt"
	"strings"
)
type n struct {
	T []string
}
func main() {
	//test := ""
	final := n{nil}
	//final := SplitPassport(test)
	fmt.Println(final)
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
		for _, p := range r { //13-25
			serial = append(serial, p[21:32])
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
