package member

import (
	"net/http"
	"reflect"

	"github.com/jarlyyn/herb-go-experimental/user-role"

	"github.com/jarlyyn/herb-go-experimental/cache-cachedmap"

	"context"

	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/cache-session"
	"github.com/herb-go/herb/user"
)

const prefixCacheBanned = "B"
const prefixCacheAccount = "A"
const prefixCacheRevoke = "V"
const prefixCacheData = "D"
const prefixCacheRole = "R"

const DefaultSessionUIDFieldName = "herb-member-uid"
const DefaultSessionRevokeFieldName = "herb-member-revoke"

type ContextType string

var DefaultContextName = ContextType("members")

type Service struct {
	SessionStore           *session.Store
	SessionUIDFieldName    string
	SessionRevokeFieldName string
	ContextName            ContextType
	Cache                  cache.Cacheable
	BannedService          BannedService
	AccountsService        AccountsService
	RevokeService          RevokeService
	PasswordService        PasswordService
	RoleService            RolesService
	DataServices           map[string]reflect.Type
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

func (s *Service) Revoke() *ServiceRevoke {
	return &ServiceRevoke{
		service: s,
	}
}

func (s *Service) Data() *ServiceData {
	return &ServiceData{
		service: s,
	}
}
func (s *Service) Roles() *ServiceRole {
	return &ServiceRole{
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
func (s *Service) RevokeField() *session.Field {
	var fieldName = s.SessionRevokeFieldName
	if fieldName == "" {
		fieldName = DefaultSessionRevokeFieldName
	}
	return s.SessionStore.Field(fieldName)
}
func (s *Service) IdentifyRequest(r *http.Request) (uid string, err error) {
	uid, err = s.UIDField().IdentifyRequest(r)
	if err != nil {
		return "", err
	}
	if uid == "" {
		return "", nil
	}
	var members = s.GetMembersFromRequest(r)
	if s.RevokeService != nil {
		_, err = members.LoadRevokeTokens(uid)
		if err != nil {
			return "", err
		}
		var revoke string
		err = s.RevokeField().Get(r, &revoke)
		if err != nil {
			return "", err
		}
		if revoke != members.RevokeTokens[uid] {
			return "", nil
		}
	}
	return uid, nil
}
func (s *Service) Logout(r *http.Request) error {
	return s.UIDField().Logout(r)
}

func (s *Service) Authorizer(rs role.RuleService) user.Authorizer {
	return &Authorizer{
		Service: s,
	}
}
func (s *Service) LoginRequiredMiddleware(unauthorizedAction http.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return user.LoginRequiredMiddleware(s, unauthorizedAction)
}

func (s *Service) LogoutMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return user.LogoutMiddleware(s)
}

func (s *Service) AuthorizeMiddleware(rs role.RuleService, unauthorizedAction http.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return user.AuthorizeMiddleware(s.Authorizer(rs), unauthorizedAction)
}

func (s *Service) BannedMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return user.AuthorizeMiddleware(s.Authorizer(nil), nil)
}

func (s *Service) RolesAuthorizeMiddleware(ruleNames ...string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var rs = role.NewRoles(ruleNames...)
	return s.AuthorizeMiddleware(rs, nil)
}

func (s *Service) Login(r *http.Request, id string) error {
	err := s.UIDField().Login(r, id)
	if err != nil {
		return err
	}
	if s.RevokeService != nil {
		member := s.GetMembersFromRequest(r)
		tokens, err := member.LoadRevokeTokens(id)
		if err != nil {
			return err
		}
		err = s.RevokeField().Set(r, tokens[id])
		if err != nil {
			return err
		}
	}
	return nil
}
