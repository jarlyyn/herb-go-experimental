package guarder

import "github.com/herb-go/herb/user/httpuser"

type Guarder interface {
	httpuser.Authorizer
	httpuser.Identifier
}
