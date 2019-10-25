package simpleapi

import (
	"net/http"

	"github.com/herb-go/fetch"
	"github.com/herb-go/util/httpserver"
)

type Server struct {
	httpserver.Config
	Name string
}
type Option struct {
	Server Server
	API    fetch.Server
	fetch.Clients
	Method  string
	Channel string
}

func (o *Option) server() *apiServer {
	return server(o.Server.Name)
}

func (o *Option) ApplyServer() error {
	if o.Server.IsEmpty() {
		return nil
	}
	return server(o.Server.Name).SetConfig(&o.Server.Config)
}

func (o *Option) Start(handler func(w http.ResponseWriter, r *http.Request)) error {
	return server(o.Server.Name).Start(o.Channel, MethodMiddleware(o.Method, handler))
}

func (o *Option) Stop() error {
	return server(o.Server.Name).Stop(o.Channel)
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
