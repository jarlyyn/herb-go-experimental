package guarder

import (
	"github.com/herb-go/fetch"
	"github.com/herb-go/herb/user/httpuser"
)

type Client struct {
	Clients            *fetch.Clients
	CredentialProvider httpuser.CredentialProvider
}

type CredentialOption struct {
	Clients fetch.Clients
	Driver  string
	Config  Option
}

type CredentialDriver interface {
	CreateProvider() (httpuser.CredentialProvider, error)
}
