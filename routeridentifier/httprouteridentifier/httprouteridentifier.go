package httprouteridentifier

import (
	"net/http"
	"net/url"

	"github.com/jarlyyn/herb-go-experimental/routeridentifier"
	"github.com/julienschmidt/httprouter"
)

type Indentifier struct {
	Enabled    bool
	Router     *httprouter.Router
	SubRouters map[string]*httprouter.Router
}

func NewIndentifier() *Indentifier {
	return &Indentifier{
		Router:     httprouter.New(),
		SubRouters: map[string]*httprouter.Router{},
	}
}

type emptWriter struct {
}

func (w *emptWriter) Header() http.Header {
	return http.Header{}
}
func (w *emptWriter) Write([]byte) (int, error) {
	return 0, nil
}
func (w *emptWriter) WriteHeader(statusCode int) {

}
func (i *Indentifier) MustIdentifyRouter(host string, r *http.Request) {
	if i.Enabled {
		var err error
		req := new(http.Request)
		req.Host = host
		req.Method = r.Method
		req.URL, err = url.Parse(r.RequestURI)
		if err != nil {
			panic(err)
		}
		i.Router.ServeHTTP(&emptWriter{}, req)
		id := routeridentifier.GetIdentificationFromRequest(req)
		if id != nil {
			routeridentifier.SetIdentificationToRequest(r, id)
		}
	}
}
