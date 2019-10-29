package guarder

import "net/http"

type Guarder interface {
	Authorize(r *http.Request) (bool, error)
	IdentifyRequest(r *http.Request) (string, error)
}
