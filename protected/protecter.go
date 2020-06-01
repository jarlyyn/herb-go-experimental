package protected

import (
	"context"
	"net/http"
)

var DefaultProtecterFunc = func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	http.Error(w, http.StatusText(404), 404)
}

type Protecter interface {
	Protect(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type ProtecterFunc func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (f ProtecterFunc) Protect(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	f(w, r, next)
}
func SetProtecter(r *http.Request, p Protecter) {
	if p == nil {
		return
	}
	ctx := context.WithValue(r.Context(), ContextKeyProtecter, p)
	*r = *r.WithContext(ctx)
}

func GetProtecter(r *http.Request) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	v := r.Context().Value(ContextKeyProtecter)
	if v != nil {
		return v.(Protecter).Protect
	}
	return DefaultProtecterFunc
}

func Protect(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	GetProtecter(r)(w, r, next)
}
