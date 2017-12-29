package member

import (
	"net/http"
	"reflect"

	"github.com/jarlyyn/herb-go-experimental/cache-cachedmap"

	"context"

	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/cache-session"
)

const prefixCacheBanned = "B"
const prefixCacheAccount = "A"
const prefixCacheSalt = "S"
const prefixCacheData = "D"
const prefixCacheRole = "R"

const DefaultSessionUIDFieldName = "herb-member-uid"
const DefaultSessionSaltFieldName = "herb-member-salt"

type ContextType string

var DefaultContextName = ContextType("members")

type Service struct {
	SessionStore         *session.Store
	SessionUIDFieldName  string
	SessionSaltFieldName string
	ContextName          ContextType
	Cache                cache.Cacheable
	BannedService        *BannedService
	AccountsService      *AccountsService
	PasswordService      *PasswordService
	RoleService          *RoleService
	DataServices         map[string]reflect.Type
}

func (s *Service) Accounts() *ServiceAccounts {
	return &ServiceAccounts{
		service: s,
	}
}
func (s *Service) Password() *ServicePassword {
	return &ServicePassword{
		service: s,
	}
}

func (s *Service) Banned() *ServiceBanned {
	return &ServiceBanned{
		service: s,
	}
}

func (s *Service) Data() *ServiceData {
	return &ServiceData{
		service: s,
	}
}

func (s *Service) RegisterData(key string, data cachedmap.CachedMap) error {
	var value = reflect.Indirect(reflect.ValueOf(data))
	if value.Kind() != reflect.Map {
		return ErrRegisteredDataNotMap
	}
	s.DataServices[key] = value.Type().Elem()
	return nil
}

func (s *Service) Middlewares() *Middlewares {
	return &Middlewares{
		service: s,
	}
}
func (s *Service) GetMembersFromRequest(r *http.Request) (members *Members) {
	var contextName = s.ContextName
	if contextName == "" {
		contextName = DefaultContextName
	}
	var membersInterface = r.Context().Value(contextName)
	if membersInterface != nil {
		if members, ok := membersInterface.(*Members); ok == true {
			return members
		}
	}
	members = NewMembers(s)
	var ctx = context.WithValue(r.Context(), contextName, members)
	*r = *r.WithContext(ctx)
	return members
}
func (s *Service) UIDField() *session.Field {
	var fieldName = s.SessionUIDFieldName
	if fieldName == "" {
		fieldName = DefaultSessionUIDFieldName
	}
	return s.SessionStore.Field(fieldName)
}
func (s *Service) SaltField() *session.Field {
	var fieldName = s.SessionSaltFieldName
	if fieldName == "" {
		fieldName = DefaultSessionSaltFieldName
	}
	return s.SessionStore.Field(fieldName)
}
func (s *Service) IdentifyRequest(r *http.Request) (uid string, err error) {
	uid, err = s.UIDField().IdentifyRequest(r)
	if err != nil {
		return "", err
	}
	var members = s.GetMembersFromRequest(r)
	if s.BannedService != nil {
		err = members.LoadBanned(uid)
		if err != nil {
			return "", err
		}
		if members.BannedMap[uid] == true {
			return "", nil
		}
	}
	if s.PasswordService != nil {
		err = members.LoadSalt(uid)
		if err != nil {
			return "", err
		}
		var salt, err = s.SaltField().IdentifyRequest(r)
		if err != nil {
			return "", err
		}
		if members.Salts[uid] == "" || salt != members.Salts[uid] {
			return "", nil
		}
	}
	return uid, nil
}
func (s *Service) Logout(r *http.Request) error {
	return s.UIDField().Logout(r)
}
