package httprouteridentifier

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"

	"github.com/herb-go/herb/middleware/router"
	"github.com/jarlyyn/herb-go-experimental/routeridentifier"
)

type Action struct {
	Enabled     bool
	Method      string
	IsSubRouter bool
	Path        string
	ID          string
	Tags        []string
}

func (a *Action) ApplyTo(i *Indentifier, router *httprouter.Router) error {
	if a.Enabled {
		if a.IsSubRouter {
			if !strings.Contains(a.Path, ":path") {
				return errors.New("httprouteridentifier: router action config must contains ':path' in path field (" + a.Path + ") if IsSUbRouter setted to true.")
			}
		}
	}
	switch a.Method {
	case "GET":
		router.GET(a.Path, a.NewAction(i))
	case "PUT":
		router.PUT(a.Path, a.NewAction(i))
	case "POST":
		router.POST(a.Path, a.NewAction(i))
	case "DELETE":
		router.DELETE(a.Path, a.NewAction(i))
	case "PATCH":
		router.PATCH(a.Path, a.NewAction(i))
	case "OPTIONS":
		router.OPTIONS(a.Path, a.NewAction(i))
	case "HEAD":
		router.HEAD(a.Path, a.NewAction(i))
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
			id := routeridentifier.GetIdentificationFromRequest(r)
			if id == nil {
				id = routeridentifier.NewIdentification()
				routeridentifier.SetIdentificationToRequest(r, id)
			}
			for k := range a.Tags {
				id.AddTag(a.Tags[k])
			}
			if !a.IsSubRouter {
				id.ID = a.MakeID(r)
				return
			}
			sr, ok := i.SubRouters[id.ID]
			if ok == false {
				return
			}
			urlparam := router.GetParams(r).Get("path")
			r.URL, err = url.Parse(urlparam)
			if err != nil {
				panic(err)
			}
			sr.ServeHTTP(w, r)
		}
	}
}

type Router []Action

func (r Router) ApplyTo(i *Indentifier, router *httprouter.Router) error {
	for _, v := range r {
		err := v.ApplyTo(i, router)
		if err != nil {
			return err
		}
	}
	return nil
}

type Config struct {
	Enabled    bool
	Router     Router
	SubRouters map[string]Router
}

func NewConfig() *Config {
	return &Config{
		SubRouters: map[string]Router{},
	}
}
func (c *Config) ApplyTo(i *Indentifier) error {
	var err error
	i.Enabled = c.Enabled
	err = c.Router.ApplyTo(i, i.Router)
	if err != nil {
		panic(err)
	}
	for k := range c.SubRouters {
		router := httprouter.New()
		err = c.SubRouters[k].ApplyTo(i, router)
		if err != nil {
			panic(err)
		}
		i.SubRouters[k] = router
	}
	return nil
}
