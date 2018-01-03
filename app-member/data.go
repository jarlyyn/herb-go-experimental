package member

import (
	"github.com/herb-go/herb/cache"
	cachedmap "github.com/jarlyyn/herb-go-experimental/cache-cachedmap"
)

type ServiceData struct {
	service *Service
}

func (s *ServiceData) Cache(field string) cache.Cacheable {
	return cache.NewNode(s.service.DataCache, field)
}
func (s *ServiceData) Clean(field string, uid string) error {
	return s.Cache(field).Del(uid)
}

func (s *ServiceData) Load(field string, data cachedmap.CachedMap, keys ...string) error {
	return cachedmap.LoadCachedMap(
		data,
		s.Cache(field),
		keys...,
	)
}
