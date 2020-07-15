package token

import (
	"time"
)

type Creator interface {
	Create(Owner, Secret, *time.Time) (*Token, error)
}

func GeneratAndCreate(c Creator, g Generator, owner Owner, expired *time.Time) (*Token, error) {
	secret, err := g.Generate()
	if err != nil {
		return nil, err
	}
	return c.Create(owner, secret, expired)
}
