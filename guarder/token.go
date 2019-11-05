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
func (g *TokenMap) Guarder() (Guarder, error) {
	return g, nil
}
func (g *TokenMap) Authorize(r *http.Request) (bool, error) {
	return IDTokenLoaderGuarderAuthorize(g, r)
}
func (g *TokenMap) IdentifyRequest(r *http.Request) (string, error) {
	return IDTokenLoaderGuarderIdentifyRequest(g, r)
}
func (g *TokenMap) Credential(id string, r *http.Request) error {
	return IDTokenLoaderGuarderCredential(g, id, r)
}

func TokenMapFactory(conf Config, prefix string) (GuarderDriver, error) {
	d := NewTokenMap()
	return d, nil
}
