package member

import (
	"github.com/herb-go/herb/cache"
	cachedmap "github.com/jarlyyn/herb-go-experimental/cache-cachedmap"
)

type Salts map[string]string
type PasswordService interface {
	Salts(uid ...string) (Salts, error)
	VerifyPassword(uid string, password string) (bool, error)
	UpdatePassword(uid string, password string) error
}
type ServicePassword struct {
	service *Service
}

func (s *ServicePassword) Cache() cache.Cacheable {
	return cache.NewNode(s.service.Cache, prefixCacheSalt)
}

func (s *ServicePassword) CleanSaltCache(uid string) error {
	return s.Cache().Del(uid)
}

func (s *ServicePassword) loader(salts Salts) func(keys ...string) error {
	return func(keys ...string) error {
		data, err := s.service.PasswordService.Salts(keys...)
		if err != nil {
			return err
		}
		for k := range data {
			salts[k] = data[k]
		}
		return nil
	}
}
func (s *ServicePassword) Load(salts Salts, keys ...string) error {
	return cachedmap.Load(
		salts,
		s.Cache(),
		s.loader(salts),
		func(key string) error {
			salts[key] = ""
			return nil
		},
		keys...,
	)
}
