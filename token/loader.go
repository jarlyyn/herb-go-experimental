package token

import (
	"bytes"
)

type Loader interface {
	Load(id ID) (*Token, error)
}

func Identify(l Loader, token *Token) (Owner, error) {
	t, err := l.Load(token.ID)
	if err != nil {
		return "", err
	}
	if t.ID == token.ID && bytes.Compare(t.Secret, token.Secret) == 0 {
		return t.Owner, nil
	}
	return "", nil
}
