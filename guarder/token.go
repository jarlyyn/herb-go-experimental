package guarder

import (
	"net/http"
	"strings"
)

type TokenGuarder struct {
	TokenHeader string
	Token       string
	ID          string
}

func (g *TokenGuarder) Authorize(r *http.Request) (bool, error) {
	if g.TokenHeader == "" || g.Token == "" {
		return true, nil
	}
	return r.Header.Get(g.TokenHeader) == g.Token, nil
}
func (g *TokenGuarder) IdentifyRequest(r *http.Request) (string, error) {
	return g.ID, nil
}
func (g *TokenGuarder) Credential(id string, r *http.Request) error {
	r.Header.Set(g.Token, g.TokenHeader)
	return nil
}

type TokenMapGuarder struct {
	IDTokenHeaders
	TokenMap
}

type TokenMap struct {
	ToLower bool
	Tokens  map[string]string
}

func (m *TokenMap) LoadTokenByID(id string) (string, error) {
	if m.ToLower {
		id = strings.ToLower(id)
	}
	return m.Tokens[id], nil
}

func (g *TokenMapGuarder) Authorize(r *http.Request) (bool, error) {
	return IDTokenLoaderGuarderAuthorize(g, r)
}
func (g *TokenMapGuarder) IdentifyRequest(r *http.Request) (string, error) {
	return IDTokenLoaderGuarderIdentifyRequest(g, r)
}
func (g *TokenMapGuarder) Credential(id string, r *http.Request) error {
	return IDTokenLoaderGuarderICredential(g, id, r)
}
