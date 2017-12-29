package member

import (
	"github.com/herb-go/herb/cache"
	cachedmap "github.com/jarlyyn/herb-go-experimental/cache-cachedmap"
)

type Banned bool
type BannedMap map[string]Banned

type BannedService struct {
	Banned func(uid ...string) (BannedMap, error)
	Ban    func(uid string, enabeld bool) error
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
			bannedMap[key] = false
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
