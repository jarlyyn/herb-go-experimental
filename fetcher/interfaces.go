package fetcher

import (
	"io"
	"net/http"
)

type RequestBuilder func(*http.Request) error

type TargetGetter interface {
	RequestMethod() string
	RequestURL() string
	RequestBody() io.Reader
	RequestBuilders() []func(*http.Request) error
}

type TargetSetter interface {
	SetRequestMethod(string)
	SetRequesetURL(string)
	SetRequestBody(io.Reader)
	SetRequestBuilders([]func(*http.Request) error)
}

type Target interface {
	TargetGetter
	TargetSetter
}

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

type Result interface {
	FromResponse(*http.Response)
	Response() *http.Response
	Error() string
}

type APIErrCode interface {
	GetAPIErrCode(err error) string
	CompareAPIErrCode(err error, code interface{}) bool
}
type ResultWithErrCode interface {
	Result
	APIErrCode
}
type TargetBuilder interface {
	BuildTarget(method string, url string, body io.Reader, builders []func(*http.Request) error) (string, string, io.Reader, []func(*http.Request) error, error)
}
