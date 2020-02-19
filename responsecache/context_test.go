package responsecache

import (
	"net/http"
	"testing"
)

func TestContext(t *testing.T) {
	r, _ := http.NewRequest("GET", "127.0.0.1", nil)
	ctx := DefaultContextField.GetContext(r)
	ctx.Identifier = func(r *http.Request) string {
		return "test"
	}
	ctx2 := DefaultContextField.GetContext(r)
	if ctx2.Identifier(r) != ctx.Identifier(r) {
		t.Fatal(ctx2)
	}
}
