package member

import (
	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/cache-session"
	"github.com/herb-go/herb/user"
)

const prefixCacheMember = "M"
const prefixCacheAccount = "A"
const prefixCacheData = "D"

type Member struct {
	UID      string
	Nickname string
	Enabled  bool
}

type MemberService struct {
	Has                func(uid string) (bool, error)
	CheckAvailablefunc func(uid string) (bool, error)
	Member             func(uid string) (Member, error)
}
type AccountsService struct {
	Accounts     func(uid string) ([]user.UserAccount, error)
	AccountToUID func(user.UserAccount) (string, error)
}
type DataService struct {
	Data func(uid string, v interface{}) error
}
type Service struct {
	UIDField        *session.Field
	ContentName     string
	Cache           cache.Cacheable
	MemberService   *MemberService
	AccountsService *AccountsService
	DataService     map[string]*DataService
}

func (s *Service) MemberCache() cache.Cacheable {
	return cache.NewNode(s.Cache, prefixCacheMember)
}

func (s *Service) CacheMember(m Member) error {
	return s.MemberCache().Set(m.UID, m, cache.DefualtTTL)
}
func (s *Service) CleanUserMemberCache(uid string) error {
	return s.MemberCache().Del(uid)
}
func (s *Service) CleanUserAccountsCache(uid string) error {
	return s.AccountsCache().Del(uid)
}
func (s *Service) GetUserMember(uid string) (*Member, error) {
	var member = Member{}
	err := s.AccountsCache().Get(uid, &member)
	if err == cache.ErrNotFound {
		member, err := s.MemberService.Member(uid)
		if err == nil {
			s.CacheMember(member)
		}
	}
	return &member, err
}
func (s *Service) AccountsCache() cache.Cacheable {
	return cache.NewNode(s.Cache, prefixCacheAccount)
}
func (s *Service) CacheUserAccounts(uid string, accounts ...user.UserAccount) error {
	return s.AccountsCache().Set(uid, accounts, cache.DefualtTTL)
}
func (s *Service) CleanUserAccountsCache(uid string) error {
	return s.AccountsCache().Del(uid)
}
func (s *Service) GetUserAccounts(uid string) ([]user.UserAccount, error) {
	var accounts = []user.UserAccount{}
	err := s.AccountsCache().Get(uid, &accounts)
	if err == cache.ErrNotFound {
		accounts, err := s.AccountsService.Accounts(uid)
		if err != nil {
			s.CacheUserAccounts(uid, accounts...)
		}
	}
	return accounts, err
}
