package concurrent

import (
	"github.com/zl14917/MastersProject/pkg/kvstore/maps"
	"sync"
)

type ThreadsafeMap struct {
	sync.RWMutex
	Map map[maps.Key]maps.Value
}

type ThreadsafeMapIterator struct {
	keys   []maps.Key
	values []maps.Value

	currentIndex int
	len          int
}

func NewThreadsafeMap() maps.Map {
	return &ThreadsafeMap{
		Map: make(map[maps.Key]maps.Value),
	}
}

func (m *ThreadsafeMap) Get(key maps.Key) (value maps.Value, ok bool) {
	m.RLock()
	defer m.RUnlock()
	value, ok = m.Map[key]
	return
}

func (m *ThreadsafeMap) Remove(key maps.Key) (value maps.Value, ok bool) {
	m.Lock()
	defer m.Unlock()
	value, ok = m.Map[key]
	delete(m.Map, key)
	return
}

func (m *ThreadsafeMap) Put(key maps.Key, value maps.Value) (err error) {
	m.Lock()
	defer m.Unlock()
	m.Map[key] = value
	return nil
}

func (m *ThreadsafeMap) Contains(key maps.Key) (ok bool) {
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

func (m *ThreadsafeMap) makeSnapshot() (keys []maps.Key, values []maps.Value) {
	m.RLock()
	defer m.RUnlock()

	len := len(m.Map)
	keys = make([]maps.Key, len, len)
	values = make([]maps.Value, len, len)

	i := 0
	for k, v := range m.Map {
		keys[i] = k
		values[i] = v
	}

	return
}

func (m *ThreadsafeMap) Iterator() maps.MapIterator {
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

func (i *ThreadsafeMapIterator) Current() (maps.Key, maps.Value) {
	return i.keys[i.currentIndex], i.values[i.currentIndex]
}
