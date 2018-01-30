package auth

import (
	"context"
	"crypto/rand"
	"errors"
	"net/http"
	"net/url"
	"strings"
)

const DefaultLoginPrefix = "/login/"
const DefaultAuthPrefix = "/auth/"

var ErrAuthParamsError = errors.New("external auth params error")
var TokenMask = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789_-")

type ContextName string

const ResultContextName = ContextName("authresult")

func DefaultNotFoundAction(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}

type Driver interface {
	ExternalLogin(provider *Provider, w http.ResponseWriter, r *http.Request)
	AuthRequest(provider *Provider, r *http.Request) (*Result, error)
}

type Auth struct {
	ProviderManager ProviderManager
	Host            string
	Path            string
	LoginPrefix     string
	AuthPrefix      string
	NotFoundAction  func(w http.ResponseWriter, r *http.Request)
	Session         Session
}

func (a *Auth) GetProviderManager() ProviderManager {
	if a.ProviderManager != nil {
		return a.ProviderManager
	}
	return DefaultProviderManager
}
func (a *Auth) RegisterProvider(keyword string, driver Driver) (*Provider, error) {
	return a.GetProviderManager().RegisterProvider(a, keyword, driver)
}
func (a *Auth) MustRegisterProvider(keyword string, driver Driver) *Provider {
	s, err := a.GetProviderManager().RegisterProvider(a, keyword, driver)
	if err != nil {
		panic(err)
	}
	return s
}
func (a *Auth) GetProvider(keyword string) (*Provider, error) {
	return a.GetProviderManager().GetProvider(a, keyword)
}
func (a *Auth) MustGetProvider(keyword string) *Provider {
	s, err := a.GetProviderManager().GetProvider(a, keyword)
	if err != nil {
		panic(err)
	}
	return s
}
func New() *Auth {
	return &Auth{}
}
func (a *Auth) Init(path string, Session Session) error {
	u, err := url.Parse(path)
	if err != nil {
		return err
	}
	*a = Auth{
		LoginPrefix:    DefaultLoginPrefix,
		AuthPrefix:     DefaultAuthPrefix,
		Session:        Session,
		Host:           u.Scheme + "://" + u.Host,
		Path:           u.Path,
		NotFoundAction: DefaultNotFoundAction,
	}
	return nil

}
func (a *Auth) MustInit(path string, Session Session) *Auth {
	err := a.Init(path, Session)
	if err != nil {
		panic(err)
	}
	return a

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
func (a *Auth) RandToken(length int) ([]byte, error) {
	token := make([]byte, length)
	_, err := rand.Read(token)
	if err != nil {
		return nil, err
	}
	l := len(TokenMask)
	for k, v := range token {
		index := int(v) % l
		token[k] = TokenMask[index]
	}
	return token, nil
}
func (a *Auth) SetResult(r *http.Request, result *Result) {
	ctx := context.WithValue(r.Context(), ResultContextName, result)
	*r = *r.WithContext(ctx)
}
func (a *Auth) Serve(SuccessAction func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var provider *Provider
		var keyword string
		path := r.URL.Path
		if keyword = strings.TrimPrefix(path, a.LoginPrefix); len(keyword) < len(path) {
			provider, err = a.GetProvider(keyword)
			if err != nil {
				panic(err)
			}
			if provider != nil {
				provider.Login(w, r)
				return
			}
		} else if keyword = strings.TrimPrefix(path, a.AuthPrefix); len(keyword) < len(path) {
			provider, err = a.GetProvider(keyword)
			if err != nil {
				panic(err)
			}
			if provider != nil {
				result, err := provider.AuthRequest(r)
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
