package responsecache

import (
	"net/http"
	"strings"
	"time"

	"github.com/herb-go/herb/cache"
)

type ParamFunc func(r *http.Request) (param string, success bool)

func (p ParamFunc) GetParams(r *http.Request) (param string, success bool) {
	return p(r)
}

type Param interface {
	GetParams(r *http.Request) (param string, success bool)
}

type FixedParam string

func (p FixedParam) GetParams(r *http.Request) (param string, success bool) {
	return string(p), true
}

var MethodParam = ParamFunc(func(r *http.Request) (param string, success bool) {
	return r.Method, true
})

type QueryParam string

func (p QueryParam) GetParams(r *http.Request) (param string, success bool) {
	return r.URL.Query().Get(string(p)), true
}

type ParamsContextBuilder struct {
	params    []Param
	TTL       time.Duration
	Validator func(ctx *Context) bool
	Cache     cache.Cacheable
}

func (b *ParamsContextBuilder) Identifier(r *http.Request) string {
	if len(b.params) == 0 {
		return ""
	}
	results := make([]string, len(b.params))
	for k := range b.params {
		s, ok := b.params[k].GetParams(r)
		if ok == false {
			return ""
		}
		results[k] = s
	}
	return strings.Join(results, cache.KeyPrefix)
}

func (b *ParamsContextBuilder) BuildContext(ctx *Context) {
	ctx.Identifier = b.Identifier
	ctx.Validator = b.Validator
	ctx.TTL = b.TTL
}
func (b *ParamsContextBuilder) Clone() *ParamsContextBuilder {
	params := make([]Param, len(b.params))
	copy(params, b.params)
	return &ParamsContextBuilder{
		params:    params,
		Validator: b.Validator,
		TTL:       b.TTL,
	}
}
func (b *ParamsContextBuilder) WithTTL(ttl time.Duration) *ParamsContextBuilder {
	pcb := b.Clone()
	pcb.TTL = ttl
	return pcb
}
func (b *ParamsContextBuilder) AppendParams(params ...Param) *ParamsContextBuilder {
	pcb := b.Clone()
	pcb.params = append(pcb.params, params...)
	return pcb
}

func (b *ParamsContextBuilder) WithValidator(v func(ctx *Context) bool) *ParamsContextBuilder {
	pcb := b.Clone()
	b.Validator = v
	return pcb
}
func NewParamsContextBuilder() *ParamsContextBuilder {
	return &ParamsContextBuilder{}
}
