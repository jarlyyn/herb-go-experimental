package responsecache

import (
	"net/http"

	"github.com/herb-go/herb/middleware"
)

type ContextField string

var DefaultContextField = "responsecache"

type Context struct {
	http.ResponseWriter
	Request     *http.Request
	ID          string
	Buffer      []byte
	BufferError error
	validated   bool
	StatusCode  int
	Validator   func(*Context) bool
}

func (c *Context) Prepare(w http.ResponseWriter, r *http.Request) {
	c.Request = r
	c.ResponseWriter = w

}

func (c *Context) WriteHeader(statusCode int) {
	c.StatusCode = statusCode
	c.ResponseWriter.WriteHeader(statusCode)
}

func (c *Context) Write(bs []byte) (int, error) {
	p, err := c.ResponseWriter.Write(bs)
	if err != nil {
		return 0, err
	}
	if c.ID != "" && c.Validator(c) {
		c.validated = true
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
