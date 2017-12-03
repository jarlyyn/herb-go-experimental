package auth

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/herb-go/herb/cache-session"
)

const DefaultLoginPrefix = "/login/"
const DefaultAuthPrefix = "/auth/"

var ErrAuthParamsError = errors.New("external auth params error")

type ContextName string

const ResultContextName = ContextName("authresult")

type DataIndex string

const DataIndexName = DataIndex("Name")
const DataIndexEmail = DataIndex("Email")
const DataIndexNickname = DataIndex("Nickname")
const DataIndexAvatar = DataIndex("Avatar")
const DataIndexProfileURL = DataIndex("ProfileURL")
const DataIndexAccessToken = DataIndex("AccessToken")
const DataIndexGender = DataIndex("ProfileURL")
const GenderMale = "M"
const GenderFemale = "F"

func DefaultNotFoundAction(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

type Driver interface {
	ExternalLogin(service *Service, w http.ResponseWriter, r *http.Request)
	AuthRequest(service *Service, r *http.Request) (*Result, error)
}

type Auth struct {
	ServiceManager ServiceManager
	Host           string
	Path           string
	LoginPrefix    string
	AuthPrefix     string
	NotFoundAction func(w http.ResponseWriter, r *http.Request)
	SessionStore   session.Store
}

func (a *Auth) GetServiceManager() ServiceManager {
	if a.ServiceManager != nil {
		return a.ServiceManager
	}
	return DefaultServiceManager
}
func (a *Auth) RegisterService(keyword string, driver Driver) (*Service, error) {
	return a.GetServiceManager().RegisterService(a, keyword, driver)
}
func (a *Auth) MustRegisterService(keyword string, driver Driver) *Service {
	s, err := a.GetServiceManager().RegisterService(a, keyword, driver)
	if err != nil {
		panic(err)
	}
	return s
}
func (a *Auth) GetService(keyword string) (*Service, error) {
	return a.GetServiceManager().GetService(a, keyword)
}
func New(path string, store session.Store) (*Auth, error) {
	u, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	a := Auth{
		LoginPrefix:    DefaultLoginPrefix,
		AuthPrefix:     DefaultAuthPrefix,
		SessionStore:   store,
		Host:           u.Scheme + "://" + u.Host,
		Path:           u.Path,
		NotFoundAction: DefaultNotFoundAction,
	}
	return &a, nil
}

func (a *Auth) MustGetResult(req *http.Request) *Result {
	data := req.Context().Value(ResultContextName)
	if data != nil {
		result, ok := data.(*Result)
		if ok {
			return result
		}
	}
	return &Result{}
}

func (a *Auth) SetResult(r *http.Request, result *Result) {
	ctx := context.WithValue(r.Context(), ResultContextName, result)
	*r = *r.WithContext(ctx)
}
func (a *Auth) Serve(SuccessAction func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var service *Service
		var keyword string
		path := r.URL.Path
		if keyword = strings.TrimPrefix(r.RequestURI, a.LoginPrefix); len(path) < len(keyword) {
			service, err = a.GetService(keyword)
			if err != nil {
				panic(err)
			}
			if service != nil {
				service.Login(w, r)
				return
			}
		} else if keyword = strings.TrimPrefix(r.RequestURI, a.AuthPrefix); len(path) < len(keyword) {
			service, err = a.GetService(keyword)
			if err != nil {
				panic(err)
			}
			if service != nil {
				result, err := service.AuthRequest(r)
				if err == ErrAuthParamsError {
					a.NotFoundAction(w, r)
					return
				}
				if err != nil {
					panic(err)
				}
				if result != nil && result.Account != "" {
					result.Keyword = keyword
					a.SetResult(r, result)
					SuccessAction(w, r)
					return
				}
			}
		}
		a.NotFoundAction(w, r)
	}
}
