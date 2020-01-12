package pattern

import (
	"net/http"
	"path/filepath"
)

type Exts map[string]bool

func (e *Exts) IsEmpty() bool {
	return len(*e) == 0
}

func (e *Exts) Match(r *http.Request) (bool, error) {
	return (*e)[filepath.Ext(r.URL.Path)], nil
}
