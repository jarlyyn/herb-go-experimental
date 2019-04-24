package routeridentifier

import "net/http"

type Identifier interface {
	MustIdentifyRouter(host string, r *http.Request)
}

var Debug bool

var DebugHeader = "herbgo-router-identification"

func Middleware(i Identifier, host string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		i.MustIdentifyRouter(host, r)
		if Debug {
			id := GetIdentificationFromRequest(r)
			if id != nil {
				w.Header().Add(DebugHeader, id.String())
			}
		}
		next(w, r)
	}
}
