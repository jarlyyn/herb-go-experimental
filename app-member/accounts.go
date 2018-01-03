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

func (s *ServiceAccounts) loader(accounts *Accounts) func(keys ...string) error {
	return func(keys ...string) error {
		data, err := s.service.AccountsService.Accounts(keys...)
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
	return s.service.AccountsService.Register(account)
}
func (s *ServiceAccounts) AccountToUID(account user.UserAccount) (uid string, err error) {
	return s.service.AccountsService.AccountToUID(account)
}
func (s *ServiceAccounts) AccountToUIDOrRegister(account user.UserAccount) (uid string, err error) {
	return s.service.AccountsService.AccountToUIDOrRegister(account)
}
func (s *ServiceAccounts) BindAccounts(uid string, account user.UserAccount) error {
	err := s.Clean(uid)
	if err != nil {
		return err
	}
	return s.service.AccountsService.BindAccounts(uid, account)
}
func (s *ServiceAccounts) UnbindAccounts(uid string, account user.UserAccount) error {
	err := s.Clean(uid)
	if err != nil {
		return err
	}
	return s.service.AccountsService.UnbindAccounts(uid, account)
}

type AccountsService interface {
	Accounts(uid ...string) (Accounts, error)
	AccountToUID(account user.UserAccount) (uid string, err error)
	Register(account user.UserAccount) (uid string, err error)
	AccountToUIDOrRegister(account user.UserAccount) (uid string, err error)
	BindAccounts(uid string, account user.UserAccount) error
	UnbindAccounts(uid string, account user.UserAccount) error
}
