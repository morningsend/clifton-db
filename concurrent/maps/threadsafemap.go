package maps

import (
	"sync"
	"time"
)

type mapValue struct {
	Data      []byte
	Delete    bool
	Timestamp int64
}

type ThreadsafeMap struct {
	sync.RWMutex
	Map map[Key]mapValue
}

type ThreadsafeMapIterator struct {
	keys   []Key
	values []Value

	currentIndex int
	len          int
}

func NewThreadsafeMap() Map {
	return &ThreadsafeMap{
		Map: make(map[Key]mapValue),
	}
}

func (m *ThreadsafeMap) Get(key Key) (value Value, ok bool) {
	m.RLock()
	defer m.RUnlock()
	readValue, ok := m.Map[key]
	if !ok || readValue.Delete {
		return nil, false
	}
	return readValue.Data, true
}

func (m *ThreadsafeMap) Remove(key Key) (value Value, ok bool) {
	m.Lock()
	defer m.Unlock()
	writeValue, ok := m.Map[key]
	if !ok {
		return nil, false
	}
	value = writeValue.Data
	writeValue.Data = nil
	writeValue.Delete = true
	writeValue.Timestamp = time.Now().Unix()
	m.Map[key] = writeValue

	return value, true
}

func (m *ThreadsafeMap) Put(key Key, value Value) (err error) {
	m.Lock()
	defer m.Unlock()
	m.Map[key] = mapValue{
		Data:   value,
		Delete: false,
	}
	return nil
}

func (m *ThreadsafeMap) Contains(key Key) (ok bool) {
	m.RLock()
	defer m.RUnlock()
	_, ok = m.Map[key]
	return false
}

func (m *ThreadsafeMap) Len() int {
	m.RLock()
	defer m.RUnlock()
	return len(m.Map)
}

func (m *ThreadsafeMap) makeSnapshot() (keys []Key, values []Value) {
	m.RLock()
	defer m.RUnlock()

	len := len(m.Map)
	keys = make([]Key, len, len)
	values = make([]Value, len, len)
	
	i := 0
	for k, v := range m.Map {
		keys[i] = k
		values[i] = v.Data
	}

	return
}

func (m *ThreadsafeMap) Iterator() MapIterator {
	keys, values := m.makeSnapshot()
	return &ThreadsafeMapIterator{
		keys:   keys,
		values: values,
		len:    len(keys),
		currentIndex: 0,
	}
}

func (i *ThreadsafeMapIterator) Next() bool {
	if i.currentIndex >= i.len {
		return false
	}

	i.currentIndex++
	return true
}

func (i *ThreadsafeMapIterator) Current() (Key, Value) {
	return i.keys[i.currentIndex], i.values[i.currentIndex]
}
