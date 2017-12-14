package member

import (
	"errors"

	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/cache-session"
	"github.com/herb-go/herb/user"
)

const prefixCacheMember = "M"
const prefixCacheAccount = "A"
const prefixCacheData = "D"
const prefixCacheCollection = "C"

var ErrProfileServiceNotFound = errors.New("error data service not found")

type Member struct {
	UID     string
	Enabled bool
}

type MemberService struct {
	Has                func(uid string) (bool, error)
	CheckAvailablefunc func(uid string) (bool, error)
	Member             func(uid ...string) ([]Member, error)
}
type AccountsService struct {
	Accounts     func(uid ...string) ([][]user.UserAccount, error)
	AccountToUID func(user.UserAccount) (string, error)
}
type ProfileService func(uid ...string) (profile []user.Profile, err error)
type Service struct {
	UIDField        *session.Field
	ContentName     string
	Cache           cache.Cacheable
	MemberService   *MemberService
	AccountsService *AccountsService
	ProfileServices map[string]ProfileService
}

func (s *Service) MemberCache() cache.Cacheable {
	return cache.NewNode(s.Cache, prefixCacheMember)
}

func (s *Service) cacheMember(m Member) error {
	return s.MemberCache().Set(m.UID, m, cache.DefualtTTL)
}
func (s *Service) CleanUserMemberCache(uid string) error {
	return s.MemberCache().Del(uid)
}
func (s *Service) GetUserMember(uid string) (*Member, error) {
	var member = Member{}
	err := s.AccountsCache().Get(uid, &member)
	if err == cache.ErrNotFound {
		members, err := s.MemberService.Member(uid)
		if err == nil {
			member = members[0]
			s.cacheMember(member)
		}
	}
	return &member, err
}
func (s *Service) AccountsCache() cache.Cacheable {
	return cache.NewNode(s.Cache, prefixCacheAccount)
}
func (s *Service) cacheUserAccounts(uid string, accounts ...user.UserAccount) error {
	return s.AccountsCache().Set(uid, accounts, cache.DefualtTTL)
}
func (s *Service) CleanUserAccountsCache(uid string) error {
	return s.AccountsCache().Del(uid)
}
func (s *Service) GetUserAccounts(uid string) ([]user.UserAccount, error) {
	var accounts = []user.UserAccount{}
	err := s.AccountsCache().Get(uid, &accounts)
	if err == cache.ErrNotFound {
		accountsList, err := s.AccountsService.Accounts(uid)
		if err != nil {
			accounts = accountsList[0]
			s.cacheUserAccounts(uid, accounts...)
		}
	}
	return accounts, err
}
func (s *Service) ProfileCache(field string) cache.Cacheable {
	c := cache.NewNode(s.Cache, prefixCacheData)
	return cache.NewNode(c, field)
}
func (s *Service) cacheProfile(field string, uid string, v interface{}) error {
	c := cache.NewNode(s.Cache, prefixCacheData)
	return cache.NewNode(c, field).Set(uid, v, 0)
}
func (s *Service) CleanProfileCache(field string, uid string) error {
	c := cache.NewNode(s.Cache, prefixCacheData)
	return cache.NewNode(c, field).Del(uid)
}

func (s *Service) UserCollection(uid string) *cache.Collection {
	c := cache.NewNode(s.Cache, prefixCacheData)
	return cache.NewCollection(c, uid, 0)
}

func (s *Service) FlushUser(uid string) error {
	err := s.MemberCache().Del(uid)
	if err != nil {
		return err
	}
	s.AccountsCache().Del(uid)
	if err != nil {
		return err
	}
	for k, _ := range s.ProfileServices {
		err = s.ProfileCache(k).Del(uid)
		if err != nil {
			return err
		}
	}
	err = s.UserCollection(uid).Flush()
	return err
}
func (s *Service) Profile(field string, uid string) (profile user.Profile, err error) {
	profile = user.Profile{}
	var c = s.ProfileCache(field)
	err = c.Get(uid, &profile)
	if err == cache.ErrNotFound {
		v, ok := s.ProfileServices[field]
		if ok == false {
			return nil, ErrProfileServiceNotFound
		}
		profileList, err := v(uid)
		if err != nil {
			return nil, err
		}
		profile = profileList[0]
		err = c.Set(uid, v, cache.DefualtTTL)
	}
	return profile, err
}

type UserResult struct {
	Member   Member
	Accounts []user.UserAccount
	Profiles map[string]user.Profile
	service  *Service
}

func newUserResult(service *Service) *UserResult {
	return &UserResult{
		Member:   Member{},
		Accounts: []user.UserAccount{},
		Profiles: map[string]user.Profile{},
		service:  service,
	}
}
func (s *Service) ListUser(uid []string, withAccounts bool, fields []string) (map[string]UserResult, error) {
	var err error
	result := map[string]UserResult{}
	uncachedMembers := []string{}
	uncachedAccounts := []string{}
	uncachedProfile := map[string][]string{}

	for _, v := range uid {
		if _, ok := result[v]; ok == true {
			continue
		}
		userresult := UserResult{
			Member:   Member{},
			Accounts: []user.UserAccount{},
			Profiles: map[string]user.Profile{},
		}
		err := s.MemberCache().Get(v, &(userresult.Member))
		if err == cache.ErrNotFound {
			err = nil
			uncachedMembers = append(uncachedMembers, v)
		}
		if err != nil {
			return nil, err
		}
		if withAccounts {
			err = s.AccountsCache().Get(v, &(userresult.Accounts))
			if err == cache.ErrNotFound {
				err = nil
				uncachedAccounts = append(uncachedAccounts, v)
			}
			if err != nil {
				return nil, err
			}
		}
		for _, field := range fields {
			profile := user.Profile{}
			if _, ok := userresult.Profiles[field]; ok == false {
				userresult.Profiles[field] = profile
			}
			err = s.ProfileCache(field).Get(v, &(profile))
			if err == cache.ErrNotFound {
				err = nil
				if _, ok := uncachedProfile[field]; ok == false {
					uncachedProfile[field] = []string{}
				}
				uncachedProfile[field] = append(uncachedProfile[field], v)
			}
			if err != nil {
				return nil, err
			}
		}
	}
}
