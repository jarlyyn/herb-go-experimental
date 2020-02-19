package responsecache

import (
	"net/http"

	"github.com/herb-go/herb/cache"
)

type Service struct {
	Cache        cache.Cacheable
	ContextField ContextField
}

func (s *Service) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ctx := s.ContextField.GetContext(r)
	if ctx.Identifier != nil && ctx.Validator != nil {
		id := ctx.Identifier(r)
		if id != "" {
			page := &cached{}
			err := s.Cache.Load(id, page, ctx.TTL, func(key string) (interface{}, error) {
				ctx.Prepare(w, r)
				next(ctx.NewWriter(), r)
				page = cacheContext(ctx)
				return page, nil
			})
			if err != nil {
				if err != cache.ErrEntryTooLarge && err != cache.ErrNotCacheable {
					panic(err)
				}
			}
			if ctx.validated {
				return
			}
			err = page.Output(w)
			if err != nil {
				panic(err)
			}
			return
		}
	}
	next(w, r)
}
