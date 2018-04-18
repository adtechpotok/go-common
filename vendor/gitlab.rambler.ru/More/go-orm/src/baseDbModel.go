package src

import (
	"sync"
)

const StatusActive = "Active"

type BaseDbModel struct {
	cache      map[int]interface{}
	cacheMutex *sync.RWMutex
	isInit     bool
}

func (m *BaseDbModel) GetCache() interface{} {
	m.init()
	return m.cache
}

func (m *BaseDbModel) init() {
	if m.isInit != true {
		m.cacheMutex = &sync.RWMutex{}
		m.cache = make(map[int]interface{})
		m.isInit = true

	}
}

func (m *BaseDbModel) FindInCache(id int) PotokOrm {
	m.init()
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()
	if _, ok := m.cache[id]; ok {
		result := m.cache[id].(PotokOrm)
		if result.IsActive() {
			return result
		}

	}

	return nil

}

func (m *BaseDbModel) FindInactiveInCache(id int) PotokOrm {
	m.init()
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()
	if _, ok := m.cache[id]; ok {
		result := m.cache[id].(PotokOrm)
		return result

	}
	return nil
}

func (m *BaseDbModel) AddToCache(v PotokOrm) {
	m.init()
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	m.cache[v.GetId()] = v

}

func (m *BaseDbModel) ClearCache() {
	m.init()
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	m.cache = make(map[int]interface{})
}

func (m *BaseDbModel) Len() int {
	return len(m.cache)
}

type PotokOrm interface {
	IsActive() bool
	GetId() int
}
