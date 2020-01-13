package requestfeature

import (
	"net/http"
)

type RequestFeature interface {
	MatchRequest(r *http.Request) (bool, error)
}
