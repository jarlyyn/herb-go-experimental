package guarder

import (
	"net/http"

	"github.com/herb-go/herb/user"
)

type IDTokenHeaders struct {
	IDHeader    string
	TokenHeader string
}

func (h *IDTokenHeaders) RequestIDTokenEnabled() (bool, error) {
	return h.IDHeader != "" && h.TokenHeader != "", nil
}
func (h *IDTokenHeaders) LoadIDTokenFormRequest(r *http.Request) (string, string, error) {
	return r.Header.Get(h.IDHeader), r.Header.Get(h.TokenHeader), nil
}
func (h *IDTokenHeaders) SetIDTokenToRequest(r *http.Request, id string, token string) error {
	r.Header.Set(h.IDHeader, id)
	r.Header.Set(h.TokenHeader, token)
	return nil
}

type RequestIDToken interface {
	RequestIDTokenEnabled() (bool, error)
	LoadIDTokenFormRequest(*http.Request) (string, string, error)
	SetIDTokenToRequest(*http.Request, string, string) error
}
type TokenLoader interface {
	LoadTokenByID(id string) (string, error)
}

type IDTokenLoaderGuarder interface {
	RequestIDToken
	TokenLoader
}

func IDTokenLoaderGuarderAuthorize(g IDTokenLoaderGuarder, r *http.Request) (bool, error) {
	e, err := g.RequestIDTokenEnabled()
	if err != nil {
		return false, err
	}
	if !e {
		return true, nil
	}
	id, token, err := g.LoadIDTokenFormRequest(r)
	if err != nil {
		return false, err
	}
	t, err := g.LoadTokenByID(id)
	return t != "" && t == token, nil

}
func IDTokenLoaderGuarderIdentifyRequest(g IDTokenLoaderGuarder, r *http.Request) (string, error) {
	id, _, err := g.LoadIDTokenFormRequest(r)
	if err != nil {
		return "", err
	}
	return id, nil
}
func IDTokenLoaderGuarderICredential(g IDTokenLoaderGuarder, id string, r *http.Request) error {
	e, err := g.RequestIDTokenEnabled()
	if err != nil {
		return err
	}
	if !e {
		return nil
	}
	if id != "" {
		token, err := g.LoadTokenByID(id)
		if err != nil {
			return err
		}
		if token != "" {
			return g.SetIDTokenToRequest(r, id, token)
		}
	}
	return user.ErrUserNotExists
}
