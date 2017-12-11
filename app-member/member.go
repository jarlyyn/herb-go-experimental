package member

import (
	"errors"
	"reflect"

	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/cache-session"
	"github.com/herb-go/herb/user"
)

const prefixCacheMember = "M"
const prefixCacheAccount = "A"
const prefixCacheData = "D"

var ErrDataServiceNotFound = errors.New("error data service not found")
var ErrDataMustBePtr = errors.New("error must post a pointer to get data")
var ErrDataNil = errors.New("error null data")
var ErrDataTypeNotMatch = errors.New("user data type not match")

type Member struct {
	UID     string
	Enabled bool
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
type DataService func(uid string) (v interface{}, err error)
type Service struct {
	UIDField        *session.Field
	ContentName     string
	Cache           cache.Cacheable
	MemberService   *MemberService
	AccountsService *AccountsService
	DataServices    map[string]DataService
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
func (s *Service) DataCache(field string) cache.Cacheable {
	c := cache.NewNode(s.Cache, prefixCacheData)
	return cache.NewNode(c, field)
}
func (s *Service) CacheData(field string, uid string, v interface{}) error {
	c := cache.NewNode(s.Cache, prefixCacheData)
	return cache.NewNode(c, field).Set(uid, v, 0)
}
func (s *Service) CleanDataCache(field string, uid string) error {
	c := cache.NewNode(s.Cache, prefixCacheData)
	return cache.NewNode(c, field).Del(uid)
}
func (s *Service) Data(field string, uid string, v interface{}) (err error) {
	var data interface{}
	var c = s.DataCache(field)
	err = c.Get(uid, v)
	if err == cache.ErrNotFound {
		ds, ok := s.DataServices[field]
		if ok == false {
			return ErrDataServiceNotFound
		}
		data, err = ds(uid)
		if err != nil {
			return err
		}
		err = cp(v, data)
		if err != nil {
			return err
		}
		err = c.Set(uid, v, cache.DefualtTTL)
	}
	return err
}

func cp(v interface{}, data interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	var vv = reflect.ValueOf(v)
	if vv.Kind() != reflect.Ptr {
		return ErrDataMustBePtr
	}
	if data == nil {
		return ErrDataNil
	}
	var d = reflect.Indirect(reflect.ValueOf(data))
	if vv.Elem().Type() != d.Type() {
		return ErrDataTypeNotMatch
	}
	vv.Elem().Set(d)
	return nil

}
