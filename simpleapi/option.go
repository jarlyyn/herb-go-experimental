package simpleapi

import (
	"net/http"

	"github.com/herb-go/fetch"
	"github.com/herb-go/util/httpserver"
)

type Option struct {
	Server httpserver.Config
	API    fetch.Server
	fetch.Clients
	Method  string
	Channel string
}

func (o *Option) Start(handler func(w http.ResponseWriter, r *http.Request)) error {
	return Start(o.Channel, MethodMiddleware(o.Method, handler))
}

func (o *Option) Stop() error {
	return Stop(o.Channel)
}

func (o *Option) EndPoint() *fetch.EndPoint {
	return o.API.EndPoint(o.Method, o.Channel)
}

type TokenOption struct {
	Option
	Token
}

func (o *TokenOption) Start(handler func(w http.ResponseWriter, r *http.Request)) error {
	return o.Option.Start(o.Token.Wrap(handler))
}

func (o *TokenOption) EndPoint() *fetch.EndPoint {
	ep := o.Option.EndPoint()
	ep.Headers.Set(o.Token.TokenHeader, o.Token.Token)
	return ep
}
