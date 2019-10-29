package guarder

import (
	"github.com/herb-go/fetch"
	"github.com/herb-go/herb/user/httpuser"
)

type Client struct {
	Clients            fetch.Clients
	CredentialProvider httpuser.CredentialProvider
}

type CredentialDriver interface {
}
