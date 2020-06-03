package protected

import (
	"net/http"

	"github.com/herb-go/herb/user/identifier/httpidentifier"
)

type Action struct {
	http.Handler
	*httpidentifier.Protecter
}

func Wrap(h http.Handler, p *httpidentifier.Protecter) *Action {
	return &Action{
		Handler:   h,
		Protecter: p,
	}
}
