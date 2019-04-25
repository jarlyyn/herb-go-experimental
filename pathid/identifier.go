package pathid

import (
	"fmt"
	"net/http"
)

type Identifier interface {
	MustIdentifyRouter(namespace string, r *http.Request)
}

var Debug bool

var DebugHeader = "herbgo-router-identification"

func Middleware(i Identifier, namespace string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if Debug {
		if namespace == "" {
			fmt.Println("Pathid served without namespace ")
		} else {
			fmt.Println("Pathid served with namespace '" + ":" + namespace + "'")
		}
	}
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		i.MustIdentifyRouter(namespace, r)
		if Debug {
			id := GetIdentificationFromRequest(r)
			if id != nil {
				w.Header().Add(DebugHeader, id.String())
			}
		}
		next(w, r)
	}
}
