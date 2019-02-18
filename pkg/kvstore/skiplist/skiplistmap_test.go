package skiplist

import "testing"

func TestSkipListMap(t *testing.T){

	m := NewSkipListMap()
	err := m.Put("hello", []byte("abcd"))
	if err != nil {
		t.Error(err)
		return
	}

	has := m.Contains("hello")

	if !has{
		t.Errorf("inserting 'hello', should contain 'hello'")
	}

}