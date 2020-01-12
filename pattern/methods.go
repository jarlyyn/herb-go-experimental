package pattern

import "net/http"

type Methods map[string]bool

func (m *Methods) IsEmpty() bool {
	return len(*m) == 0
}

func (m *Methods) Match(r *http.Request) (bool, error) {
	return (*m)[r.Method], nil
}
