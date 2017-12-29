package member

import (
	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/user"
	cachedmap "github.com/jarlyyn/herb-go-experimental/cache-cachedmap"
)

type Accounts map[string]user.UserAccounts

type ServiceAccounts struct {
	service *Service
}

func (s *ServiceAccounts) Loader(accounts Accounts) func(keys ...string) error {
	return func(keys ...string) error {
		data, err := s.service.AccountsService.Accounts(keys...)
		if err != nil {
			return err
		}
		for k := range data {
			accounts[k] = data[k]
		}
		return nil
	}
}

func (s *ServiceAccounts) Cache() cache.Cacheable {
	return cache.NewNode(s.service.Cache, prefixCacheAccount)
}

func (s *ServiceAccounts) Clean(uid string) error {
	return s.Cache().Del(uid)
}
func (s *ServiceAccounts) Load(accounts Accounts, keys ...string) error {
	return cachedmap.Load(
		accounts,
		s.Cache(),
		s.Loader(accounts),
		func(key string) error {
			accounts[key] = user.UserAccounts{}
			return nil
		},
		keys...,
	)
}

type AccountsService struct {
	Accounts               func(uid ...string) (Accounts, error)
	AccountToUID           func(user.UserAccount) (uid string, err error)
	Register               func(accounts Accounts) (uid string, err error)
	AccountToUIDOrRegister func(user.UserAccount) (uid string, err error)
	BindAccounts           func(uid string, accounts Accounts) error
	UnbindAccounts         func(uid string, accounts Accounts) error
}
