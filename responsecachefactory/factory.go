package responsecachefactory

import (
	"fmt"
	"time"

	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/middleware"
	"github.com/herb-go/herb/middleware/middlewarefactory"
	"github.com/herb-go/responsecache"
)

type Config struct {
	Params      []string
	Validator   *string
	Cache       string
	TTLInSecond int64
}

func New() *Factory {
	return &Factory{
		RegisteredParams: map[string]responsecache.Param{},
		Validators:       map[string]func(ctx *responsecache.Context) bool{},
		Caches:           map[string]cache.Cacheable{},
	}
}

type Factory struct {
	RegisteredParams map[string]responsecache.Param
	Validators       map[string]func(ctx *responsecache.Context) bool
	Caches           map[string]cache.Cacheable
}

func (f *Factory) WithParam(name string, p responsecache.Param) *Factory {
	f.RegisteredParams[name] = p
	return f
}

func (f *Factory) WithValidator(name string, validator func(ctx *responsecache.Context) bool) *Factory {
	f.Validators[name] = validator
	return f
}

func (f *Factory) WithCache(name string, c cache.Cacheable) *Factory {
	f.Caches[name] = c
	return f
}

func (f *Factory) CreateMiddleware(name string, loader func(v interface{}) error) (middleware.Middleware, error) {
	c := &Config{}
	params := make([]responsecache.Param, len(c.Params))
	for k, v := range c.Params {
		p := f.RegisteredParams[v]
		if p == nil {
			return nil, fmt.Errorf("response cache factory: %s param not registered", p)
		}
		params[k] = p
	}
	pb := responsecache.NewParamsContextBuilder()

	pb.AppendParams(responsecache.FixedParam(name))
	pb.AppendParams(params...)
	if c.Validator != nil {
		v := f.Validators[*c.Validator]
		if v == nil {
			return nil, fmt.Errorf("response cache factory: %s validator not registered", *c.Validator)
		}
		pb.WithValidator(v)
	}
	pb.WithTTL(time.Duration(c.TTLInSecond) * time.Second)
	cacheable := f.Caches[c.Cache]
	if cacheable == nil {
		return nil, fmt.Errorf("response cache factory: %s validator not registered", c.Cache)
	}
	pb.WithCache(cacheable)
	return responsecache.New(pb), nil
}

func (f *Factory) Factory() middlewarefactory.Factory {
	return middlewarefactory.FactoryFunc(f.CreateMiddleware)
}
