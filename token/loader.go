package token

import (
	"bytes"
)

type Loader interface {
	Load(id ID) (*Token, error)
}

func Identify(l Loader, id ID, secret Secret) (Owner, error) {
	t, err := l.Load(id)
	if err != nil {
		if err == ErrIDNotFound {
			return "", nil
		}
		return "", err
	}
	if t.ID == id && bytes.Compare(t.Secret, secret) == 0 {
		return t.Owner, nil
	}
	return "", nil
}
