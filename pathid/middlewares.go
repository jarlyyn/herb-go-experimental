package pathid

import (
	"net/http"
)

type Middlewares struct {
	Enabled           bool
	RouterMiddlewares map[string][]func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	TagMiddlewares    map[string][]func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	IDMiddlewares     map[string][]func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}
