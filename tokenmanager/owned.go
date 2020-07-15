package tokenmanager

import (
	"time"

	"github.com/herb-go/herb/user/token"
)

type OwnedToken interface {
	Owner() string
	Token() *token.Token
}

type TemporaryToken interface {
	ExpiredAt() *time.Time
	Token() *token.Token
}
