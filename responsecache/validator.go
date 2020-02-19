package responsecache

var DefaultValidator = func(ctx *Context) bool {
	return ctx.StatusCode < 500
}
