package auth

import "net/http"

type Service struct {
	Driver  Driver
	Auth    *Auth
	Keyword string
}

func (s *Service) Login(w http.ResponseWriter, r *http.Request) {
	s.Driver.ExternalLogin(s, w, r)
}

func (s *Service) AuthRequest(r *http.Request) (*Result, error) {
	return s.Driver.AuthRequest(s, r)
}

func (s *Service) LoginUrl() string {
	return s.Auth.Host + s.Auth.Path + s.Auth.LoginPrefix + s.Keyword
}

func (s *Service) AuthUrl() string {
	return s.Auth.Host + s.Auth.Path + s.Auth.AuthPrefix + s.Keyword
}
