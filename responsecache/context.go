package responsecache

import (
	"context"
	"net/http"
	"time"

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

var DefaultContextField = ContextField("responsecache")

type Context struct {
	http.ResponseWriter
	Request    *http.Request
	ID         string
	TTL        time.Duration
	Buffer     []byte
	validated  bool
	StatusCode int
	Validator  func(*Context) bool
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
	if c.ID != "" && c.Validator(c) {
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
