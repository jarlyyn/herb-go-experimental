package tokenmanager

import (
	"github.com/herb-go/herb/user/token"
)

type Manager struct {
	Generator token.Generator
	Loader    token.Loader
	Encoding  token.Encoding
	Revoker   token.Revoker
}
