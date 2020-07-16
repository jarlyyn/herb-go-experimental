package token

import "time"

type Store interface {
	Loader
	Refresh(ID, *time.Time) error
	Insert(*Token, *time.Time) error
	Update(ID, Secret) error
	Revoke(id ID) error
}

type Manager interface {
	Store
	Creator
	Regenerator
}

func CreateManagedToken(m Manager, o Owner, expiredat *time.Time) (*Token, error) {
	t, err := m.Create(o)
	if err != nil {
		return nil, err
	}
	err = m.Regenerate(t)
	if err != nil {
		return nil, err
	}
	err = m.Insert(t, expiredat)
	if err != nil {
		return nil, err
	}
	return t, nil
}
