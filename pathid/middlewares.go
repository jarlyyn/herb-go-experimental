package pathid

import (
	"net/http"
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

func (r *MiddlewaresRouter) MustRegister(rule Rule, f MiddlewaresFactory) {
	if !rule.Enabled || rule.ID == "" {
		return
	}
	var middlewaresMap map[string][]func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	switch rule.Type {
	case RuleTypeParent:
		middlewaresMap = r.ParentMiddlewares
	case RuleTypeTag:
		middlewaresMap = r.TagMiddlewares
	case RuleTypeID, "":
		middlewaresMap = r.IDMiddlewares
	default:
		return
	}
	middlewareList, ok := middlewaresMap[rule.ID]
	if !ok {
		middlewareList = []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){}
	}
	middlewareList = append(middlewareList, f.MustCreateMiddlewares()...)
	middlewaresMap[rule.ID] = middlewareList
}

type Rule struct {
	Enabled bool
	Type    string
	ID      string
}
