package responsecache

import (
	"context"
	"net/http"
	"time"

	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/middleware"
)

type ContextField string

func (c ContextField) GetContext(r *http.Request) *Context {
	var ctx *Context
	v := r.Context().Value(c)
	if v == nil {
		ctx = NewContext()
		reqctx := context.WithValue(r.Context(), c, ctx)
		req := r.WithContext(reqctx)
		*r = *req
	} else {
		ctx = v.(*Context)
	}
	return ctx
}

func (c ContextField) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := c.GetContext(r)
	if ctx.Identifier != nil && ctx.Validator != nil && ctx.Cache != nil {
		id := ctx.Identifier(r)
		if id != "" {
			page := &cached{}
			err := ctx.Cache.Load(id, page, ctx.TTL, func(key string) (interface{}, error) {
				ctx.Prepare(w, r)
				next(ctx.NewWriter(), r)
				page = cacheContext(ctx)
				return page, nil
			})
			if err != nil {
				if err != cache.ErrEntryTooLarge && err != cache.ErrNotCacheable {
					panic(err)
				}
			}
			if ctx.validated {
				return
			}
			err = page.Output(w)
			if err != nil {
				panic(err)
			}
			return
		}
	}
	next(w, r)
}

func (c ContextField) NewResponseCache(b ContextBuilder) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		ctx := c.GetContext(r)
		b.BuildContext(ctx)
		c.ServeMiddleware(w, r, next)
	}
}

func New(b ContextBuilder) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return DefaultContextField.NewResponseCache(b)
}

var DefaultContextField = ContextField("responsecache")

type Context struct {
	http.ResponseWriter
	Request    *http.Request
	Identifier func(*http.Request) string
	TTL        time.Duration
	Buffer     []byte
	validated  bool
	StatusCode int
	Validator  func(*Context) bool
	Cache      cache.Cacheable
}

func NewContext() *Context {
	return &Context{}

}
func (c *Context) Prepare(w http.ResponseWriter, r *http.Request) {
	c.Request = r
	c.ResponseWriter = w

}

func (c *Context) WriteHeader(statusCode int) {
	c.StatusCode = statusCode
	if (c.Validator != nil && c.Validator(c)) || DefaultValidator(c) {
		c.validated = true
	}
	c.ResponseWriter.WriteHeader(statusCode)
}

func (c *Context) Write(bs []byte) (int, error) {
	p, err := c.ResponseWriter.Write(bs)
	if err != nil {
		return 0, err
	}
	if c.validated {
		c.Buffer = append(c.Buffer, bs...)
	}
	return p, nil
}
func (c *Context) NewWriter() http.ResponseWriter {
	rw := middleware.WrapResponseWriter(c.ResponseWriter)
	rw.Functions().WriteFunc = c.Write
	rw.Functions().WriteHeaderFunc = c.WriteHeader
	return rw
}
