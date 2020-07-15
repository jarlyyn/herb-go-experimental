package fieldprotecter

import (
	"fmt"

	"github.com/herb-go/herb/user/credential"
	"github.com/herb-go/herb/user/credential/protecter"
	"github.com/herb-go/httpinfomanager"
	actionoverseer "github.com/herb-go/providers/herb/overseers/actionoverseer"
	"github.com/herb-go/worker"
)

type Config struct {
	Fields         map[string]*httpinfomanager.FieldName
	OnFailWorkerID string
	Verifier       string
	VerifierConfig func(v interface{}) error
}

func (c *Config) ApplyTo(p protecter.Protecter) error {
	var credentialers []protecter.Credentialer
	for k, v := range c.Fields {

		credentialers = append(credentialers, &Credentialer{credentialType: credential.Type(k), field: v.Field()})
	}
	p.Credentialers = credentialers
	if c.OnFailWorkerID != "" {
		a := actionoverseer.GetActionByID(c.OnFailWorkerID)
		if a == nil {
			return fmt.Errorf("%w (%s)", worker.ErrWorkerNotFound, c.OnFailWorkerID)
		}
		p.OnFail = a
	}
	return nil
}
