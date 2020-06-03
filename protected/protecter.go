package protected

import (
	"context"
	"net/http"

	"github.com/herb-go/herb/user/identifier/httpidentifier"
)

var DefaultProtecterFunc = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	http.Error(w, http.StatusText(404), 404)
}

func SetProtecter(r *http.Request, p *httpidentifier.Protecter) {
	if p == nil {
		return
	}
	ctx := context.WithValue(r.Context(), ContextKeyProtecter, p)
	*r = *r.WithContext(ctx)
}

func getProtecter(r *http.Request) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	v := r.Context().Value(ContextKeyProtecter)
	if v != nil {
		return v.(*httpidentifier.Protecter).ServeMiddleware
	}
	return DefaultProtecterFunc
}

func Protect(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	getProtecter(r)(w, r, next)
}
