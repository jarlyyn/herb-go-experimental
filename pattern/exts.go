package pattern

import "net/http"

type Exts map[string]bool

func (e *Exts) IsEmpty() bool {
	return len(*e) == 0
}

func (e *Exts) Match(r *http.Request) (bool, error) {
	return (*e)[r.Method], nil
}
