package auth

import "net/http"

type Provider struct {
	Driver  Driver
	Auth    *Auth
	Keyword string
}

func (p *Provider) Login(w http.ResponseWriter, r *http.Request) {
	p.Driver.ExternalLogin(p, w, r)
}

func (p *Provider) AuthRequest(r *http.Request) (*Result, error) {
	return p.Driver.AuthRequest(p, r)
}

func (p *Provider) LoginUrl() string {
	return p.Auth.Host + p.Auth.Path + p.Auth.LoginPrefix + p.Keyword
}

func (p *Provider) AuthUrl() string {
	return p.Auth.Host + p.Auth.Path + p.Auth.AuthPrefix + p.Keyword
}
