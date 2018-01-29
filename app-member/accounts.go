package member

import (
	"github.com/herb-go/herb/cache"
	cachedmap "github.com/herb-go/herb/cache-cachedmap"
	"github.com/herb-go/herb/user"
)

type Accounts map[string]user.UserAccounts

type ServiceAccounts struct {
	service *Service
}

func (s *ServiceAccounts) loader(accounts *Accounts) func(keys ...string) error {
	return func(keys ...string) error {
		data, err := s.service.AccountsProvider.Accounts(keys...)
		if err != nil {
			return err
		}
		for k := range data {
			(*accounts)[k] = data[k]
		}
		return nil
	}
}

func (s *ServiceAccounts) Cache() cache.Cacheable {
	return s.service.AccountsCache
}

func (s *ServiceAccounts) Clean(uid string) error {
	return s.Cache().Del(uid)
}
func (s *ServiceAccounts) Load(accounts *Accounts, keys ...string) error {
	return cachedmap.Load(
		accounts,
		s.Cache(),
		s.loader(accounts),
		func(key string) error {
			(*accounts)[key] = user.UserAccounts{}
			return nil
		},
		keys...,
	)
}

func (s *ServiceAccounts) Register(account user.UserAccount) (uid string, err error) {
	return s.service.AccountsProvider.Register(account)
}
func (s *ServiceAccounts) AccountToUID(account user.UserAccount) (uid string, err error) {
	return s.service.AccountsProvider.AccountToUID(account)
}
func (s *ServiceAccounts) AccountToUIDOrRegister(account user.UserAccount) (uid string, err error) {
	return s.service.AccountsProvider.AccountToUIDOrRegister(account)
}
func (s *ServiceAccounts) BindAccounts(uid string, account user.UserAccount) error {
	err := s.Clean(uid)
	if err != nil {
		return err
	}
	return s.service.AccountsProvider.BindAccounts(uid, account)
}
func (s *ServiceAccounts) UnbindAccounts(uid string, account user.UserAccount) error {
	err := s.Clean(uid)
	if err != nil {
		return err
	}
	return s.service.AccountsProvider.UnbindAccounts(uid, account)
}

type AccountsProvider interface {
	Accounts(uid ...string) (Accounts, error)
	AccountToUID(account user.UserAccount) (uid string, err error)
	Register(account user.UserAccount) (uid string, err error)
	AccountToUIDOrRegister(account user.UserAccount) (uid string, err error)
	BindAccounts(uid string, account user.UserAccount) error
	UnbindAccounts(uid string, account user.UserAccount) error
}
