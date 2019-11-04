package guarder

import "github.com/herb-go/herb/user/httpuser"

type Guarder interface {
	httpuser.Authorizer
	httpuser.Identifier
}

type GuarderDriver interface {
	Guarder() (Guarder, error)
}
type GuarderProvider struct {
	Driver GuarderDriver
}

func (g *GuarderProvider) Guarder() (Guarder, error) {
	return g.Driver.Guarder()
}
