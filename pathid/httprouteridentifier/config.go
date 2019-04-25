package httprouteridentifier

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/jarlyyn/herb-go-experimental/pathid"
)

type Action struct {
	Enabled     bool
	Method      string
	IsSubRouter bool
	Path        string
	ID          string
	Tags        []string
}

func NewAction() *Action {
	return &Action{
		Tags: []string{},
	}
}
func (a *Action) ApplyTo(i *Indentifier, router *httprouter.Router) error {
	if a.Enabled {
		if a.IsSubRouter {
			if !strings.Contains(a.Path, ":path") {
				return errors.New("httprouteridentifier: router action config must contains ':path' in path field (" + a.Path + ") if IsSUbRouter setted to true.")
			}
		}
	}
	a.Method = strings.ToUpper(a.Method)
	switch a.Method {
	case "GET", "PUT", "POST", "DELETE", "PATCH", "OPTIONS", "HEAD":
		router.Handle(a.Method, a.Path, a.NewAction(i))
	default:
		if !(a.Method == "ALL" || a.IsSubRouter) {
			return errors.New("httprouteridentifier:unkown method '" + a.Method + "' for action path '" + a.Path + "'")
		}
		router.GET(a.Path, a.NewAction(i))
		router.PUT(a.Path, a.NewAction(i))
		router.POST(a.Path, a.NewAction(i))
		router.DELETE(a.Path, a.NewAction(i))
		router.PATCH(a.Path, a.NewAction(i))
		router.OPTIONS(a.Path, a.NewAction(i))
		router.HEAD(a.Path, a.NewAction(i))
	}
	return nil
}
func (a *Action) MakeID(r *http.Request) string {
	return r.Host + "/" + a.ID + "#" + r.Method
}
func (a *Action) NewAction(i *Indentifier) func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		if a.Enabled {
			var err error
			id := pathid.GetIdentificationFromRequest(r)
			if id == nil {
				id = pathid.NewIdentification()
				pathid.SetIdentificationToRequest(r, id)
			}
			for k := range a.Tags {
				id.AddTag(a.Tags[k])
			}
			if !a.IsSubRouter {
				id.ID = a.MakeID(r)
				return
			}
			id.AddParent(a.ID)
			sr, ok := i.SubRouters[a.ID]
			if ok == false {
				return
			}
			urlparam := p.ByName("path")
			r.URL, err = url.Parse("/" + urlparam)
			if err != nil {
				panic(err)
			}
			sr.ServeHTTP(w, r)
		}
	}
}

type Router []*Action

func (r Router) ApplyTo(i *Indentifier, router *httprouter.Router) error {
	for _, v := range r {
		err := v.ApplyTo(i, router)
		if err != nil {
			return err
		}
	}
	return nil
}

func NewRouter() Router {
	return []*Action{}
}

type Config struct {
	Enabled    bool
	Router     Router
	SubRouters map[string]Router
}

func NewConfig() *Config {
	return &Config{
		Router:     []*Action{},
		SubRouters: map[string]Router{},
	}
}
func (c *Config) ApplyTo(i *Indentifier) error {
	var err error
	i.Enabled = c.Enabled
	err = c.Router.ApplyTo(i, i.Router)
	if err != nil {
		return err
	}
	for k := range c.SubRouters {
		router := httprouter.New()
		err = c.SubRouters[k].ApplyTo(i, router)
		if err != nil {
			return err
		}
		i.SubRouters[k] = router
	}
	return nil
}
