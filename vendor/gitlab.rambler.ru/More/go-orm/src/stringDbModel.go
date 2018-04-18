package src

import (
	"sync"
)

type StringDbModel struct {
	cache      map[string]PotokStringOrm
	cacheMutex *sync.RWMutex
	isInit     bool
}

func (m *StringDbModel) init() {
	if m.isInit != true {
		m.cacheMutex = &sync.RWMutex{}
		m.cache = make(map[string]PotokStringOrm)
		m.isInit = true
	}
}

func (m *StringDbModel) GetCache() interface{} {
	m.init()
	return m.cache
}

func (m *StringDbModel) FindInCache(id string) PotokStringOrm {
	m.init()
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()
	if _, ok := m.cache[id]; ok {
		result := m.cache[id].(PotokStringOrm)
		if result.IsActive() {
			return result
		}

	}

	return nil

}

func (m *StringDbModel) AddToCache(v PotokStringOrm) {
	m.init()
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	m.cache[v.GetId()] = v

}

func (m *StringDbModel) ClearCache() {
	m.init()
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	m.cache = make(map[string]PotokStringOrm)
}

func (m *StringDbModel) Len() int {
	return len(m.cache)
}

func (m *StringDbModel) ClearOutDated() {
	m.init()
	m.cacheMutex.Lock()
	defer 	m.cacheMutex.Unlock()
	data := m.cache

	for key, item := range data {
		if item.OutDated() {
			delete(data, key)
		}
	}

	m.cache = data

}

type PotokStringOrm interface {
	IsActive() bool
	GetId() string
	OutDated() bool
}
