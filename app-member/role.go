package member

import (
	"github.com/herb-go/herb/cache"
	cachedmap "github.com/jarlyyn/herb-go-experimental/cache-cachedmap"
	role "github.com/jarlyyn/herb-go-experimental/user-role"
)

type Roles map[string]role.Roles

type ServiceRole struct {
	service *Service
}
type RolesService interface {
	Roles(uid ...string) (Roles, error)
}

func (s *ServiceRole) Load(roles Roles, keys ...string) error {
	return cachedmap.Load(
		roles,
		s.Cache(),
		s.loader(roles),
		func(key string) error {
			roles[key] = role.Roles{}
			return nil
		},
		keys...,
	)
}

func (s *ServiceRole) Cache() cache.Cacheable {
	return cache.NewNode(s.service.Cache, prefixCacheRole)
}

func (s *ServiceRole) Clean(uid string) error {
	return s.Cache().Del(uid)
}
func (s *ServiceRole) loader(roles Roles) func(keys ...string) error {
	return func(keys ...string) error {
		data, err := s.service.RoleService.Roles(keys...)
		if err != nil {
			return err
		}
		for k := range data {
			roles[k] = data[k]
		}
		return nil
	}
}
