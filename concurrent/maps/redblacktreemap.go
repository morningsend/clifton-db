package maps

type RedBlackTreeMap struct {
}

func (t *RedBlackTreeMap) Remove(key Key) (value Value, ok bool) {
	panic("implement me")
}

func NewRedBlackTreeMap() Map {
	return &RedBlackTreeMap{

	}
}

func (t *RedBlackTreeMap) Get(key Key) (value Value, ok bool) {
	return nil, false
}

func (t *RedBlackTreeMap) Put(key Key, value Value) (err error) {
	return
}

func (t *RedBlackTreeMap) Contains(key Key) (ok bool) {
	return false
}

func (t *RedBlackTreeMap) Iterator() MapIterator {
	return nil
}

func (t *RedBlackTreeMap) Len() int {
	return 0
}
