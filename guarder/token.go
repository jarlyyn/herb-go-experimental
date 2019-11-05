package guarder

import (
	"net/http"
	"strings"
)

type Token struct {
	TokenHeader string
	Token       string
	ID          string
}

func (g *Token) Authorize(r *http.Request) (bool, error) {
	if g.TokenHeader == "" || g.Token == "" {
		return true, nil
	}
	return r.Header.Get(g.TokenHeader) == g.Token, nil
}
func (g *Token) IdentifyRequest(r *http.Request) (string, error) {
	return g.ID, nil
}
func (g *Token) Credential(id string, r *http.Request) error {
	r.Header.Set(g.Token, g.TokenHeader)
	return nil
}

func NewTokenMap() *TokenMap {
	return &TokenMap{}
}

type TokenMap struct {
	IDTokenHeaders
	TokenMapConfig
}

type TokenMapConfig struct {
	ToLower bool
	Tokens  map[string]string
}

func (m *TokenMapConfig) LoadTokenByID(id string) (string, error) {
	if m.ToLower {
		id = strings.ToLower(id)
	}
	return m.Tokens[id], nil
}

func (g *TokenMap) Authorize(r *http.Request) (bool, error) {
	return IDTokenLoaderGuarderAuthorize(g, r)
}
func (g *TokenMap) IdentifyRequest(r *http.Request) (string, error) {
	return IDTokenLoaderGuarderIdentifyRequest(g, r)
}

func TokenMapFactory(conf Config, prefix string) (Guarder, error) {
	var err error
	d := NewTokenMap()
	err = conf.Get("IDHeader", &d.IDHeader)
	if err != nil {
		return nil, err
	}
	err = conf.Get("TokenHeader", &d.TokenHeader)
	if err != nil {
		return nil, err
	}
	err = conf.Get("ToLower", &d.ToLower)
	if err != nil {
		return nil, err
	}
	err = conf.Get("Tokens", &d.Tokens)
	if err != nil {
		return nil, err
	}
	return d, nil
}
