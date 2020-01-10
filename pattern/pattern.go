package pattern

import (
	"net/http"
)

type Pattern interface {
	Match(r *http.Request) (bool, error)
	IsEmpty() bool
}
