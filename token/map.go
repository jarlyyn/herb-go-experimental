package token

import (
	"strconv"
	"sync"
	"time"
)

type MapData struct {
	ID        ID
	Owner     Owner
	Secret    Secret
	ExpiredAt *time.Time
}

func (d *MapData) Expired() bool {
	if d.ExpiredAt != NeverExpired && d.ExpiredAt.Before(time.Now()) {
		return true
	}
	return false
}

type Map struct {
	lastID      int64
	IDGenerator func() (string, error)
	data        map[ID]*MapData
	locker      sync.Mutex
}

func (m *Map) GenerateID() (string, error) {
	m.lastID++
	return strconv.FormatInt(m.lastID, 10), nil
}
func (m *Map) Load(id ID) (*Token, error) {
	m.locker.Lock()
	defer m.locker.Unlock()
	data, ok := m.data[id]
	if !ok || data.Expired() {
		return nil, ErrIDNotFound
	}
	token := New()
	token.Owner = data.Owner
	token.ID = data.ID
	token.Secret = data.Secret
	return token, nil
}

func (m *Map) Update(id ID, secret Secret) error {
	m.locker.Lock()
	defer m.locker.Unlock()
	data, ok := m.data[id]
	if !ok || data.Expired() {
		return ErrIDNotFound
	}
	data.Secret = secret
	return nil
}
func (m *Map) Refresh(id ID, expiredAt *time.Time) error {
	m.locker.Lock()
	defer m.locker.Unlock()
	data, ok := m.data[id]
	if !ok || data.Expired() {
		return ErrIDNotFound
	}
	data.ExpiredAt = expiredAt
	return nil
}

func (m *Map) Revoke(id ID) error {
	m.locker.Lock()
	defer m.locker.Unlock()

	delete(m.data, id)
	return nil
}
func (m *Map) Create(owner Owner, secret Secret, expiredat *time.Time) (*Token, error) {
	m.locker.Lock()
	defer m.locker.Unlock()
	idstr, err := m.GenerateID()
	if err != nil {
		return nil, err
	}
	id := ID(idstr)
	m.data[id] = &MapData{
		Owner:     owner,
		ID:        id,
		Secret:    secret,
		ExpiredAt: expiredat,
	}
	t := New()
	t.Owner = owner
	t.ID = id
	t.Secret = secret
	return t, nil
}
func (m *Map) Reset() {
	m.data = map[ID]*MapData{}
	m.lastID = 0
	m.IDGenerator = m.GenerateID
}
func NewMap() *Map {
	m := &Map{}
	m.Reset()
	return m
}
