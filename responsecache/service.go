package responsecache

import cache "github.com/patrickmn/go-cache"

type Service struct {
	Cache        cache.Cache
	ContextField ContextField
}
