package src

import (
	"sync"
)

type StringSlicedDbModel struct {
	cache      map[string][]PotokStringSliceOrm
	cacheMutex *sync.RWMutex
	isInit     bool
	lenght     int
}

func (m *StringSlicedDbModel) init() {
	if m.isInit != true {
		m.cacheMutex = &sync.RWMutex{}
		m.cache = make(map[string][]PotokStringSliceOrm)
		m.isInit = true
		m.lenght = 0

	}
}
func (m *StringSlicedDbModel) GetCache() interface{} {
	m.init()
	return m.cache
}
func (m *StringSlicedDbModel) FindInCache(id string) []PotokStringSliceOrm {
	m.init()
	m.cacheMutex.Lock()
	var result []PotokStringSliceOrm
	var t []PotokStringSliceOrm
	defer m.cacheMutex.Unlock()
	if _, ok := m.cache[id]; ok {
		result = m.cache[id]
	}

	for _, val := range result {
		if val.IsActive() {
			t = append(t, val)
		}
	}

	return t
}

func (m *StringSlicedDbModel) AddToCache(v PotokStringSliceOrm) {
	m.init()
	var res []PotokStringSliceOrm
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	current := m.cache[v.GetId()]
	m.lenght = m.lenght - len(current)
	for _, val := range current {
		if val.GetRelationKey() != v.GetRelationKey() {
			res = append(res, val)
		}
	}
	m.lenght = len(res) + 1
	m.cache[v.GetId()] = append(res, v)
}

func (m *StringSlicedDbModel) ClearCache() {
	m.init()
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	m.cache = make(map[string][]PotokStringSliceOrm)
}

func (m *StringSlicedDbModel) Len() int {
	return m.lenght
}

type PotokStringSliceOrm interface {
	GetRelationKey() int
	GetId() string
	IsActive() bool
}
