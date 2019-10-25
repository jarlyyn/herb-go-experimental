package simpleapi

import "net/http"

type Token struct {
	TokenHeader string
	Token       string
}

func (t *Token) Wrap(handler func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if t.TokenHeader != "" || t.Token != "" {
			if r.Header.Get(t.TokenHeader) != t.Token {
				http.Error(w, http.StatusText(401), 401)
				return
			}
		}
		handler(w, r)
	}
}

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
