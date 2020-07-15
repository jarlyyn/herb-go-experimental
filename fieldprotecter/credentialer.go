package fieldprotecter

import (
	"net/http"

	"github.com/herb-go/herb/middleware/httpinfo"
	"github.com/herb-go/herb/user/credential"
)

type Credentialer struct {
	credentialType credential.Type
	field          httpinfo.Field
}

func (c *Credentialer) CredentialRequest(r *http.Request) credential.Credential {
	return &Credential{
		request:      r,
		credentialer: c,
	}
}
