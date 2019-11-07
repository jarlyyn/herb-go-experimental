package guarder

import "net/http"

var DefaultOnFail = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(401), 401)
})

type Guarder struct {
	Reader     RequestParamsReader
	Identifier RequestParamsIdentifier
	OnFail     http.Handler
}

func (g *Guarder) IdentifyRequest(r *http.Request) (string, error) {
	p, err := g.Reader.ReadParamsFromRequest(r)
	if err != nil {
		return "", err
	}
	return g.Identifier.IdentifyRequestParams(p)
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

type RequestParamsGuarderOption interface {
	RequestParamsReaderDriver() string
	RequestParamsIdentifierDriver() string
	DriverConfig() *Config
}

type DriverConfig struct {
	RequestParamsIdentifierDriverField
	RequestParamsReaderDriverField
}

type RequestParamsGuarderConfigMap struct {
	DriverConfig
	Config ConfigMap
}
