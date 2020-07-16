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
	lastID int64
	Creator
	Regenerator
	data   map[ID]*MapData
	locker sync.Mutex
}

func (m *Map) GenerateID(o Owner) (*Token, error) {
	m.lastID++
	t := o.NewToken()
	t.ID = ID(strconv.FormatInt(m.lastID, 10))
	return t, nil
}

func (m *Map) Load(id ID) (*Token, error) {
	m.locker.Lock()
	defer m.locker.Unlock()
	data, ok := m.data[id]
	if !ok || data.Expired() {
		return nil, ErrTokenNotFound
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
		return ErrTokenNotFound
	}
	data.Secret = secret
	return nil
}
func (m *Map) Refresh(id ID, expiredAt *time.Time) error {
	m.locker.Lock()
	defer m.locker.Unlock()
	data, ok := m.data[id]
	if !ok || data.Expired() {
		return ErrTokenNotFound
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
func (m *Map) Insert(t *Token, expiredat *time.Time) error {
	m.locker.Lock()
	defer m.locker.Unlock()
	if t.ID == "" {
		return ErrEmptyID
	}
	_, ok := m.data[t.ID]
	if ok {
		return ErrTokenExists
	}
	m.data[t.ID] = &MapData{
		Owner:     t.Owner,
		ID:        t.ID,
		Secret:    t.Secret,
		ExpiredAt: expiredat,
	}
	return nil
}
func (m *Map) Reset() {
	m.data = map[ID]*MapData{}
	m.lastID = 0
	m.Regenerator = NopRegenerator
	m.Creator = CreatorFunc(m.GenerateID)
}
func NewMap() *Map {
	m := &Map{}
	m.Reset()
	return m
}
