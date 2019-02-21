package kvstore

type MapKVStore map[string][]byte

func NewMapKVStore() KVStore {
	return &MapKVStore{}
}

func (s MapKVStore) Get(key string) ([]byte, bool, error) {
	data, ok := s[key]
	return data, ok, nil
}

func (s MapKVStore) Put(key string, data []byte) (error) {
	s[key] = data
	return nil
}

func (s MapKVStore) Exists(key string) (bool, error) {
	_, ok := s[key]
	return ok, nil
}
