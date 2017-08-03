package transport

import (
	"sync"
)

// should be remove after go 1.9
// replaced by sync.Map

type Map struct {
	m    map[interface{}]interface{}
	lock sync.Mutex
}

func NewMap() *Map {
	return &Map{
		m: make(map[interface{}]interface{}),
	}
}

func (m *Map) Store(key, value interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.m[key] = value
}

func (m *Map) Load(key interface{}) (value interface{}, ok bool) {
	m.lock.Lock()
	defer m.lock.Unlock()
	value, ok = m.m[key]
	return
}
