package member

import (
	"github.com/herb-go/herb/cache"
	cachedmap "github.com/jarlyyn/herb-go-experimental/cache-cachedmap"
)

type BannedMap map[string]bool

type BannedService interface {
	Banned(uid ...string) (BannedMap, error)
	Ban(uid string, banned bool) error
}
type ServiceBanned struct {
	service *Service
}

func (s *ServiceBanned) Load(bannedMap BannedMap, keys ...string) error {
	return cachedmap.Load(
		bannedMap,
		s.Cache(),
		s.loader(bannedMap),
		func(key string) error {
			var val = false
			bannedMap[key] = val
			return nil
		},
		keys...,
	)
}

func (s *ServiceBanned) Cache() cache.Cacheable {
	return cache.NewNode(s.service.Cache, prefixCacheBanned)
}

func (s *ServiceBanned) Clean(uid string) error {
	return s.Cache().Del(uid)
}
func (s *ServiceBanned) Ban(uid string, banned bool) error {
	err := s.Clean(uid)
	if err != nil {
		return err
	}
	return s.service.BannedService.Ban(uid, banned)
}

func (s *ServiceBanned) loader(bannedMap BannedMap) func(keys ...string) error {
	return func(keys ...string) error {
		data, err := s.service.BannedService.Banned(keys...)
		if err != nil {
			return err
		}
		for k := range data {
			bannedMap[k] = data[k]
		}
		return nil
	}
}
