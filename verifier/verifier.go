package verifier

import (
	"time"

	"github.com/herb-go/herb/user/token"
	"github.com/herb-go/member"
)

var MemberVerifer struct {
	Member   *member.Service
	UIDLogin bool
}

var HashVerifer struct {
	Fields         map[string]string
	Required       []string
	TTL            time.Duration
	TokenFieldName string
	TokenLoader    token.Loader
	Hash           func([]byte) []byte
}

var TokenVerifier struct {
	TokenLoader token.Loader
}
