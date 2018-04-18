package src

import (
	"sync"
)

type BaseSlicedDbModel struct {
	cache      map[int][]PotokSliceOrm
	cacheMutex *sync.RWMutex
	isInit     bool
	lenght     int
}

func (m *BaseSlicedDbModel) init() {
	if m.isInit != true {
		m.cacheMutex = &sync.RWMutex{}
		m.cache = make(map[int][]PotokSliceOrm)
		m.isInit = true
		m.lenght = 0

	}
}

func (m *BaseSlicedDbModel) GetCache() interface{} {
	m.init()
	return m.cache
}

func (m *BaseSlicedDbModel) FindInCache(id int) []PotokSliceOrm {
	m.init()
	m.cacheMutex.Lock()
	var result []PotokSliceOrm
	var t []PotokSliceOrm
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

func (m *BaseSlicedDbModel) AddToCache(v PotokSliceOrm) {
	m.init()
	var res []PotokSliceOrm
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	current := m.cache[v.GetId()]
	m.lenght = m.lenght - len(current)
	for _, val := range current {
		if val.GetRelationKey() != v.GetRelationKey() {
			res = append(res, val)
		}
	}
	m.lenght = m.lenght + len(res) + 1
	m.cache[v.GetId()] = append(res, v)
}

func (m *BaseSlicedDbModel) ClearCache() {
	m.init()
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	m.cache = make(map[int][]PotokSliceOrm)
}

func (m *BaseSlicedDbModel) Len() int {
	return m.lenght
}

type PotokSliceOrm interface {
	GetRelationKey() int
	GetId() int
	IsActive() bool
}
