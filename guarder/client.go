package guarder

import (
	"net/http"
)

type Credential struct {
	Credential RequestParamsCredential
	Writer     RequestParamsWriter
}

func (c *Credential) CredentialRequest(r *http.Request) error {
	p, err := c.Credential.CredentialRequestParams()
	if err != nil {
		return err
	}
	return c.Writer.WriterParamsToRequest(r, p)
}
