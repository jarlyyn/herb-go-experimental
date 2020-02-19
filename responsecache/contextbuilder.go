package responsecache

import (
	"net/http"
	"time"
)

type ContextBuilder interface {
	BuildContext(*Context)
}

type PlainContextBuilder struct {
	Identifier func(*http.Request) string
	Validator  func(ctx *Context) bool
	TTL        time.Duration
}

func (b *PlainContextBuilder) BuildContext(ctx *Context) {
	ctx.Identifier = b.Identifier
	ctx.Validator = b.Validator
	ctx.TTL = b.TTL
}
