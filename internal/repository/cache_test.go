package repository

import "testing"


func TestRamCacheAddAndGet(t *testing.T) {
	testCases := []string{
		"test",
		"    test     ",
		"   AnotHEr T e s T    ",
	}
	cache := NewRamCache()
	for _, tC := range testCases {
		cache.AddValue(tC)
		t.Run(tC, func(t *testing.T) {
			ok := cache.GetValue(tC)
			if !ok {
				t.Error("Ошибка поиска ключа")
			}
		})
	}
}
