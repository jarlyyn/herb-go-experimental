package member

import (
	"encoding/json"
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
const prefixCacheToken = "T"
const prefixCacheData = "D"
const prefixCacheRole = "R"

const DefaultSessionUIDFieldName = "herb-member-uid"
const DefaultSessionMemberTokenFieldName = "herb-member-token"

type ContextType string

var DefaultContextName = ContextType("members")

type Service struct {
	SessionStore           *session.Store
	SessionUIDFieldName    string
	SessionMemberFieldName string
	ContextName            ContextType
	BannedService          BannedService
	BannedCache            cache.Cacheable
	AccountsService        AccountsService
	AccountsCache          cache.Cacheable
	TokenService           TokenService
	TokenCache             cache.Cacheable
	PasswordService        PasswordService
	RoleService            RolesService
	RoleCache              cache.Cacheable
	DataServices           map[string]reflect.Type
	DataCache              cache.Cacheable
	AccountTypes           map[string]user.AccountType
}

type Installable interface {
	InstallToService(service *Service)
}

func (s *Service) RegisterAccountType(keyword string, t user.AccountType) {
	s.AccountTypes[keyword] = t
}
func (s *Service) Install(i Installable) {
	i.InstallToService(s)
}
func (s *Service) NewAccount(keyword string, account string) (*user.UserAccount, error) {
	accountType, ok := s.AccountTypes[keyword]
	if ok == false {
		return nil, ErrAccountKeywordNotRegistered
	}
	return accountType.NewAccount(keyword, account)
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

func (s *Service) Token() *ServiceToken {
	return &ServiceToken{
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
	if s.DataServices == nil {
		s.DataServices = map[string]reflect.Type{}
	}
	var value = reflect.ValueOf(data)
	if value.Kind() != reflect.Map {
		return ErrRegisteredDataNotMap
	}
	s.DataServices[key] = value.Type()
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
func (s *Service) MemberTokenField() *session.Field {
	var fieldName = s.SessionMemberFieldName
	if fieldName == "" {
		fieldName = DefaultSessionMemberTokenFieldName
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
	if s.TokenService != nil {
		_, err = members.LoadTokens(uid)
		if err != nil {
			return "", err
		}
		var token string
		err = s.MemberTokenField().Get(r, &token)
		if err != nil {
			return "", err
		}
		if token != members.Tokens[uid] {
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
		Service:     s,
		RuleService: rs,
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
	if s.TokenService != nil {
		member := s.GetMembersFromRequest(r)
		tokens, err := member.LoadTokens(id)
		if err != nil {
			return err
		}
		err = s.MemberTokenField().Set(r, tokens[id])
		if err != nil {
			return err
		}
	}
	return nil
}

var dummyCache = cache.New()

func New(store *session.Store) *Service {
	return &Service{
		SessionStore:  store,
		BannedCache:   dummyCache,
		AccountsCache: dummyCache,
		TokenCache:    dummyCache,
		RoleCache:     dummyCache,
		DataCache:     dummyCache,
		DataServices:  map[string]reflect.Type{},
		AccountTypes:  map[string]user.AccountType{},
	}
}

func NewWithSubCache(store *session.Store, c cache.Cacheable) *Service {
	var s = New(store)
	s.BannedCache = cache.NewNode(c, prefixCacheBanned)
	s.AccountsCache = cache.NewNode(c, prefixCacheAccount)
	s.TokenCache = cache.NewNode(c, prefixCacheToken)
	s.RoleCache = cache.NewNode(c, prefixCacheRole)
	s.DataCache = cache.NewNode(c, prefixCacheData)
	return s
}

func init() {
	var err = dummyCache.Open("dummycache", json.RawMessage(""), 1)
	if err != nil {
		panic(err)
	}
}
