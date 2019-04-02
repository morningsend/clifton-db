package vector

import "sync"

type LockVector struct {
	lock sync.RWMutex
	data []interface{}
}

func (v *LockVector) Get(idx int) (value interface{}, ok bool) {
	v.lock.RLock()
	defer v.lock.RUnlock()

	l := len(v.data)
	if idx < 0 || idx >= l {
		return
	}
	value = v.data[idx]
	ok = true
	return
}

func (v *LockVector) Set(idx int, element interface{}) {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.data[idx] = element

}

func (v *LockVector) Append(value interface{}) {
	v.lock.Lock()
	defer v.lock.Unlock()
	v.data = append(v.data, value)
}

func (v *LockVector) Pop() (value interface{}, ok bool) {
	v.lock.Lock()
	defer v.lock.Unlock()

	length := len(v.data)
	if length < 1 {
		return nil, false
	}
	value = v.data[length-1]
	v.data = v.data[:length-1]
	return value, true
}

func (v *LockVector) Len() int {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return len(v.data)
}

func (v *LockVector) Cap() int {
	v.lock.RLock()
	defer v.lock.RUnlock()
	return cap(v.data)
}
func (v *LockVector) Resize(newLength int) {
	v.lock.Lock()
	defer v.lock.Unlock()
	if len(v.data) == newLength {
		return
	} else if len(v.data) > newLength {
		v.data = v.data[:newLength]
	} else {
		newData := make([]interface{}, newLength, newLength*2+1)
		for i, d := range v.data {
			newData[i] = d
		}

		v.data = newData
	}
}
