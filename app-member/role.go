package member

import (
	"github.com/herb-go/herb/cache"
	cachedmap "github.com/jarlyyn/herb-go-experimental/cache-cachedmap"
)

type UserRoles []string

type Roles map[string]UserRoles

type RoleService struct {
	Role func(uid ...string) (Roles, error)
}
type ServiceRole struct {
	service *Service
}

func (s *ServiceRole) Load(roles Roles, keys ...string) error {
	return cachedmap.Load(
		roles,
		s.Cache(),
		s.loader(roles),
		func(key string) error {
			roles[key] = UserRoles{}
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
		data, err := s.service.RoleService.Role(keys...)
		if err != nil {
			return err
		}
		for k := range data {
			roles[k] = data[k]
		}
		return nil
	}
}
