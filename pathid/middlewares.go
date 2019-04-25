package pathid

import (
	"net/http"

	"github.com/herb-go/herb/middleware"
)

const RuleTypeID = "id"
const RuleTypeParent = "parent"
const RuleTypeTag = "tag"

type MiddlewaresFactory interface {
	MustCreateMiddlewares() []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}
type MiddlewaresRouter struct {
	Enabled           bool
	ParentMiddlewares map[string][]func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	TagMiddlewares    map[string][]func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	IDMiddlewares     map[string][]func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

func (r *MiddlewaresRouter) getMiddlewareMapByType(t string) map[string][]func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	switch t {
	case RuleTypeParent:
		return r.ParentMiddlewares
	case RuleTypeTag:
		return r.TagMiddlewares
	case RuleTypeID, "":
		return r.IDMiddlewares
	default:
		return nil
	}
}
func (r *MiddlewaresRouter) MustRegister(rule Rule, f MiddlewaresFactory) {
	if !rule.Enabled || rule.ID == "" {
		return
	}
	middlewaresMap := r.getMiddlewareMapByType(rule.Type)
	if middlewaresMap == nil {
		return
	}
	middlewareList, ok := middlewaresMap[rule.ID]
	if !ok {
		middlewareList = []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){}
	}
	middlewareList = append(middlewareList, f.MustCreateMiddlewares()...)
	middlewaresMap[rule.ID] = middlewareList
}
func (r *MiddlewaresRouter) loadMiddlewares(id *Identification) []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var result = []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){}
	if id.ID != "" {
		middlewareList, ok := r.IDMiddlewares[id.ID]
		if ok {
			result = append(result, middlewareList...)
		}
	}
	if len(r.TagMiddlewares) > 0 {
		for k := range id.Tags {
			middlewareList, ok := r.TagMiddlewares[id.Tags[k]]
			if ok {
				result = append(result, middlewareList...)
			}
		}
	}
	if len(r.ParentMiddlewares) > 0 {
		for k := range id.Parents {
			middlewareList, ok := r.ParentMiddlewares[id.Parents[k]]
			if ok {
				result = append(result, middlewareList...)
			}
		}
	}
	return result
}
func (r *MiddlewaresRouter) ServerMiddleware(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if r.Enabled {
		id := GetIdentificationFromRequest(req)
		if id != nil {
			middleares := r.loadMiddlewares(id)
			middleware.New(middleares...).ServeMiddleware(w, req, next)
			return
		}
	}
	next(w, req)
}

type Rule struct {
	Enabled bool
	Type    string
	ID      string
}
