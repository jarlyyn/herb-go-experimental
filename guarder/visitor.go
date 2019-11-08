package guarder

import (
	"net/http"
)

type Visitor struct {
	Credential RequestParamsCredential
	Mapper     RequestParamsMapper
}

func (c *Visitor) CredentialRequest(r *http.Request) error {
	p, err := c.Credential.CredentialRequestParams()
	if err != nil {
		return err
	}
	return c.Mapper.WriteParamsToRequest(r, p)
}
