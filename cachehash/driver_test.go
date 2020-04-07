package cachehash

import (
	"encoding/binary"
	"sync"
	"testing"
)

func TestInt64(t *testing.T) {
	var v int64
	v = v - 1
	var bytes = make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(v))
	if int64(binary.BigEndian.Uint64(bytes)) != -1 {
		t.Fatal(bytes)
	}
}

type testStore struct {
	data map[string]*Hash
}

func (s *testStore) Open() error {
	return nil
}
func (s *testStore) Close() error {
	return nil
}
func (s *testStore) Flush() error {
	s.data = map[string]*Hash{}
	return nil
}
func (s *testStore) Hash(key string) (string, error) {
	return key[:1], nil
}
func (s *testStore) Lock(string) (func(), error) {
	l := sync.Mutex{}
	l.Lock()
	return func() {
		l.Unlock()
	}, nil
}
func (s *testStore) Load(hash string) (*Hash, error) {
	d := s.data[hash]
	if d == nil {
		return NewHash(), nil
	}
	return d, nil
}
func (s *testStore) Delete(hash string) error {
	delete(s.data, hash)
	return nil
}
func (s *testStore) Save(hash string, status *Status, data *Hash) error {
	s.data[hash] = data
	return nil
}

func newTestStore() Store {
	return &testStore{
		data: map[string]*Hash{},
	}
}

func newTestDriver() *Driver {
	return NewDriver(newTestStore())
}
