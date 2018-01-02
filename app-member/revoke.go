package member

import (
	"github.com/herb-go/herb/cache"
	cachedmap "github.com/jarlyyn/herb-go-experimental/cache-cachedmap"
)

type RevokeTokens map[string]string

type ServiceRevoke struct {
	service *Service
}
type RevokeService interface {
	RevokeTokens(uid ...string) (RevokeTokens, error)
	Revoke(uid string) (string, error)
}

func (s *ServiceRevoke) Cache() cache.Cacheable {
	return cache.NewNode(s.service.Cache, prefixCacheRevoke)
}

func (s *ServiceRevoke) Clean(uid string) error {
	return s.Cache().Del(uid)
}

func (s *ServiceRevoke) Revoke(uid string) (string, error) {
	err := s.Clean(uid)
	if err != nil {
		return "", err
	}
	return s.service.RevokeService.Revoke(uid)
}

func (s *ServiceRevoke) loader(revokeTokens RevokeTokens) func(keys ...string) error {
	return func(keys ...string) error {
		data, err := s.service.RevokeService.RevokeTokens(keys...)
		if err != nil {
			return err
		}
		for k := range data {
			revokeTokens[k] = data[k]
		}
		return nil
	}
}
func (s *ServiceRevoke) Load(revokeTokens RevokeTokens, keys ...string) error {
	return cachedmap.Load(
		revokeTokens,
		s.Cache(),
		s.loader(revokeTokens),
		func(key string) error {
			var val = ""
			revokeTokens[key] = val
			return nil
		},
		keys...,
	)
}
