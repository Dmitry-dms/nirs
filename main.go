package main

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"log"
	"regexp"
	"strings"
)

type n struct {
	T []string
}

func main() {
	test := "hello"

	final := generateHash(test)
	fmt.Println(final)

}
func generateHash(s string) []byte {
	g := fnv.New32()
	g.Write([]byte(s))
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, g.Sum32())
    return bs
}
func show(s []string) {
	for i, j := range s {
		fmt.Printf("[%d] %s \n", i, j)
	}
}
func splitAddress(s string) []string {
	var a []string
	addresses := strings.Split(s, ";")
	a = append(a, addresses...)
	return a
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

func split(s string) []string {
	var result []string
	pattern := "\\(+.*"
	r := regexp.MustCompile(pattern)
	secPart := r.FindAllString(s, 5) // Делим строку на две части с помощью regex, отделяем первую часть и вторую
	if secPart == nil {              // Если деление не сработало, значит имеем всего одно имя
		result = append(result, s)
		return result
	}
	secPart2 := secPart[0][1:(len(secPart[0]) - 1)] // Избавляемся от лишних круглых скобок
	splitted := strings.Split(secPart2, ";")
	for _, name := range splitted {
		k := strings.Split(name, "(") // Делим на подстроки
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
	return removeDuplicateStr(result) // Избавляемся от дубликатов имени
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
		var singlePass []string
		var serial []string
		raw := strings.Split(s[:len(s)-1], ",")
		singlePass = append(singlePass, raw...)
		for _, p := range singlePass {
			if len(p) < 9 { // Если разделение строк пошло не по шаблону, прекращаем разделение
				log.Printf("Skip password split. Unusual data: %s", raw)
				return Passport{
					RawData:      raw,
					SerialAndNum: nil,
				}
			}
			runeString := []rune(p)
			if runeString[8] == 'Р' {
				serial = append(serial, string(runeString[12:23]))
			} else if runeString[8] == 'С' {
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
func trimSpace(s []string) []string {
	var n []string
	for _, j := range s {
		m := strings.TrimSpace(j)
		n = append(n, m)
	}
	return n
}
