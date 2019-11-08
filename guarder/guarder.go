package guarder

import "net/http"

var DefaultOnFail = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(401), 401)
})

type Guarder struct {
	Mapper     Mapper
	Identifier Identifier
	OnFail     http.Handler
}

func (g *Guarder) IdentifyRequest(r *http.Request) (string, error) {
	p, err := g.Mapper.ReadParamsFromRequest(r)
	if err != nil {
		return "", err
	}
	return g.Identifier.IdentifyParams(p)
}

func (g *Guarder) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	id, err := g.IdentifyRequest(r)
	if err != nil {
		panic(err)
	}
	if id != "" {
		next(w, r)
		return
	}
	g.OnFail.ServeHTTP(w, r)
}
func NewGuarder() *Guarder {
	return &Guarder{
		OnFail: DefaultOnFail,
	}
}
