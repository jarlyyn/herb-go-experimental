package member

import (
	"github.com/herb-go/herb/cache"
	cachedmap "github.com/herb-go/herb/cache-cachedmap"
)

type Tokens map[string]string

type ServiceToken struct {
	service *Service
}
type TokenProvider interface {
	Tokens(uid ...string) (Tokens, error)
	Revoke(uid string) (string, error)
}

func (s *ServiceToken) Cache() cache.Cacheable {
	return s.service.TokenCache
}

func (s *ServiceToken) Clean(uid string) error {
	return s.Cache().Del(uid)
}

func (s *ServiceToken) Revoke(uid string) (string, error) {
	err := s.Clean(uid)
	if err != nil {
		return "", err
	}
	return s.service.TokenProvider.Revoke(uid)
}

func (s *ServiceToken) loader(Tokens *Tokens) func(keys ...string) error {
	return func(keys ...string) error {
		data, err := s.service.TokenProvider.Tokens(keys...)
		if err != nil {
			return err
		}
		for k := range data {
			(*Tokens)[k] = data[k]
		}
		return nil
	}
}
func (s *ServiceToken) Load(Tokens *Tokens, keys ...string) error {
	return cachedmap.Load(
		Tokens,
		s.Cache(),
		s.loader(Tokens),
		func(key string) error {
			(*Tokens)[key] = ""
			return nil
		},
		keys...,
	)
}
