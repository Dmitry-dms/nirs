package repository

import "strings"

type RamCache struct {
	m map[string]struct{}
}

func NewRamCache() *RamCache {
	c := RamCache{
		m: make(map[string]struct{}, 50000),
	}
	return &c
}

func (c *RamCache) AddValue(key string) {
	trimSpaceKey := strings.ReplaceAll(key, " ", "")
	lower :=strings.ToLower(trimSpaceKey)
	c.m[lower] = struct{}{}
}
func (c *RamCache) GetValue(key string) bool {
	if key == "" {
		return false
	}
	trimSpaceKey := strings.ReplaceAll(key, " ", "")
	lower :=strings.ToLower(trimSpaceKey)
	_, ok := c.m[lower]
	return ok
}
