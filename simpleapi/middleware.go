package simpleapi

import "net/http"

func MethodMiddleware(method string, handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if method != "" {
			if r.Method != method {
				http.Error(w, http.StatusText(405), 405)
				return
			}
		}
		handler(w, r)
	}
}
