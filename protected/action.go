package protected

import (
	"net/http"
)

type Action interface {
	http.Handler
	Protecter
}

type WrappedAction struct {
	http.Handler
	Protecter
}

func Wrap(h http.Handler, p Protecter) *WrappedAction {
	return &WrappedAction{
		Handler:   h,
		Protecter: p,
	}
}
