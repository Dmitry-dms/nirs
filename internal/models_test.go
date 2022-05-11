package internal

import (
	"reflect"
	"testing"
)

// Все тестовые данные являются фальшивыми


func TestTrimSuffixAndPrefix(t *testing.T) {
	testCases := []struct {
		src string
		res string
	}{
		{
			src: " <![CDATA[Яблоко ]]>",
			res: "Яблоко",
		},
		{
			src: " <![CDATA[1958 ]]>",
			res: "1958",
		},
		{
			src: " <![CDATA[]]>",
			res: "NULL",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.src, func(t *testing.T) {
			r := trimCDATASuffixAndPrefix(tC.src)
			if r != tC.res {
				t.Errorf("Неправильный формат строки: должно быть \"%s\" , имеется %s", tC.res, r)
			}
		})
	}
}

func TestSplitPassport(t *testing.T) {
	testCases := []struct {
		src          string
		serialAndNum []string
	}{
		{
			src:          "ПАСПОРТ РФ: 0000 111111 ВЫДАН ОУФМС РОССИИ --------------------------,",
			serialAndNum: []string{"0000 111111"},
		},
		{
			src:          "ПАСПОРТ РФ: 1212 333444 ВЫДАН ОВД ---------------------,",
			serialAndNum: []string{"1212 333444"},
		},
		{
			src:          "ПАСПОРТ РФ: 5555 666666 ВЫДАН -----------,ПАСПОРТ РФ: 6767 564455 ---------------,ПАСПОРТ РФ: 4444 555676 ВЫДАН ------------------,",
			serialAndNum: []string{"5555 666666", "6767 564455", "4444 555676"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.src, func(t *testing.T) {
			res := splitPassport(tC.src, nil)
			if !reflect.DeepEqual(res.SerialAndNum, tC.serialAndNum) {
				t.Errorf("Ошибка парсинга пасспорта: имеется %s необходимо %s", res.SerialAndNum, tC.serialAndNum)
			}
		})
	}
}

func TestSplitAddress(t *testing.T) {
	testCases := []struct {
		src     string
		address []string
	}{
		{
			src:     "РФ ул. Кирова кв. 133,",
			address: []string{"РФ ул. Кирова кв. 133"},
		},
		{
			src:     "РФ, Москва кв. 4; РФ,  КВ. 3",
			address: []string{"РФ, Москва кв. 4", "РФ,  КВ. 3"},
		},
		{
			src:     "РФ, Москва кв. 4; РФ,  КВ. 3; США г. Вашингтон",
			address: []string{"РФ, Москва кв. 4", "РФ,  КВ. 3", "США г. Вашингтон"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.src, func(t *testing.T) {
			res := splitAddress(tC.src)
			if !reflect.DeepEqual(res, tC.address) {
				t.Errorf("Ошибка парсинга адреса: имеется %s необходимо %s", res, tC.address)
			}
		})
	}
}

func TestSplitNames(t *testing.T) {
	testCases := []struct {
		src   string
		names []string
	}{
		{
			src:   "Иванов Иван Иванович*",
			names: []string{"Иванов Иван Иванович*"},
		},
		{
			src:   "ВКП* (Иванов И.А.; Иванов Г.А.; Иванов А.А.)",
			names: []string{"ВКП*", "Иванов И.А.", "Иванов Г.А.", "Иванов А.А."},
		},
		{
			src:   "Иванов Иван Иванович* (Иванов Дмитрий Иванович)",
			names: []string{"Иванов Иван Иванович*", "Иванов Дмитрий Иванович"},
		},
		{
			src:   "Очень плохая компания (ОПК)",
			names: []string{"Очень плохая компания", "ОПК"},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.src, func(t *testing.T) {
			res := splitNames(tC.src)
			if !isEqualWithoudOrder(res, tC.names) {
				t.Errorf("Ошибка парсинга имен: имеется %s необходимо %s", res, tC.names)
			}
		})
	}
}

func isEqualWithoudOrder[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for _, v := range a {
		founded := false
		for _, j := range b {
			if v == j {
				founded = true
			}
		}
		if !founded {
			return false
		}
	}
	return true
}

func TestRemoveDublicates(t *testing.T) {
	testCases := []struct {
		desc []string
		res  []string
	}{
		{
			desc: []string{"Иван", "Иванов", "Иван", "Володя"},
			res:  []string{"Иван", "Иванов", "Володя"},
		},
		{
			desc: []string{"Дмитрий", "Иванов", "Иван", "Володя"},
			res:  []string{"Дмитрий","Иван", "Иванов", "Володя"},
		},
		{
			desc: []string{"Дмитрий", "Дмитрий"},
			res:  []string{"Дмитрий"},
		},
	}
	for _, tC := range testCases {
		t.Run("tC.desc[i]", func(t *testing.T) {
			r := removeDuplicateStr(tC.desc)
			if !isEqualWithoudOrder(r, tC.res) {
				t.Errorf("Ошибка парсинга имен: имеется %s необходимо %s", r, tC.res)
			}
		})
	}
}
