package auth

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/herb-go/herb/cache-session"
)

const DefaultLoginPrefix = "/login/"
const DefaultAuthPrefix = "/auth/"

type ContextName string

const ResultContextName = ContextName("herbgo-thirtpartauth-result")

type DataIndex string

const DataIndexEmail = DataIndex("Email")
const DataIndexNickname = DataIndex("Nickname")
const DataIndexAvatar = DataIndex("Avatar")
const DataIndexProfileURL = DataIndex("ProfileURL")
const DataIndexAccessToken = DataIndex("AccessToken")

func DefaultNotFoundAction(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

type Driver interface {
	LoginAction() func(w http.ResponseWriter, r *http.Request)
	Auth(r *http.Request) (*Result, error)
}
type Service struct {
	Driver  Driver
	Auth    *Auth
	Keyword string
}

type Data map[DataIndex][]string

func (d *Data) Value(index DataIndex) string {
	data, ok := (*d)[index]
	if ok == false || len(data) == 0 {
		return ""
	}
	return data[0]
}

func (d *Data) Values(index DataIndex) []string {
	data, ok := (*d)[index]
	if ok == false {
		return nil
	}
	return data
}

func (d *Data) SetValue(index DataIndex, value string) {
	(*d)[index] = []string{value}
}

func (d *Data) SetValues(index DataIndex, values []string) {
	(*d)[index] = values
}

func (d *Data) AddValue(index DataIndex, value string) {
	data, ok := (*d)[index]
	if ok == false {
		data = []string{}
	}
	data = append(data, value)
	(*d)[index] = data
}

type Result struct {
	Keyword string
	Account string
	Data    Data
}

type Auth struct {
	Services       map[string]*Service
	Host           string
	Path           string
	LoginPrefix    string
	AuthPrefix     string
	NotFoundAction func(w http.ResponseWriter, r *http.Request)
	SessionStore   session.Store
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
func (a *Auth) RegisterService(keyword string, driver Driver) *Service {
	s := &Service{
		Driver:  driver,
		Auth:    a,
		Keyword: keyword,
	}
	a.Services[keyword] = s
	return s
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

func (a *Auth) GetService(keyword string) *Service {
	s, ok := a.Services[keyword]
	if ok {
		return s
	}
	return nil
}
func (a *Auth) SetResult(r *http.Request, result *Result) {
	ctx := context.WithValue(r.Context(), ResultContextName, result)
	*r = *r.WithContext(ctx)
}
func (a *Auth) Serve(SuccessAction func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var service *Service
		var keyword string
		path := r.URL.Path
		if keyword = strings.TrimPrefix(r.RequestURI, a.LoginPrefix); len(path) < len(keyword) {
			if service = a.GetService(keyword); service != nil {
				service.Driver.LoginAction()(w, r)
				return
			}
		} else if keyword = strings.TrimPrefix(r.RequestURI, a.AuthPrefix); len(path) < len(keyword) {
			if service = a.GetService(keyword); service != nil {
				result, err := service.Driver.Auth(r)
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
