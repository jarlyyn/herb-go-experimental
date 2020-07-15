package fieldprotecter

import (
	"net/http"

	"github.com/herb-go/herb/user/credential"
)

type Credential struct {
	request      *http.Request
	credentialer *Credentialer
}

func (c *Credential) Type() credential.Type {
	return c.credentialer.credentialType
}
func (c *Credential) Data() ([]byte, error) {
	data, _, err := c.credentialer.field.LoadInfo(c.request)
	return data, err
}
