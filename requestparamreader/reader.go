package requestparamreader

import "net/http"

type Reader func(r *http.Request) ([]byte, error)

type ReaderFactory interface {
	CreateReader(loader func(v interface{}) error) (Reader, error)
}

type FactoryFunc func(loader func(v interface{}) error) (Reader, error)

func (f FactoryFunc) CreateFactory(loader func(v interface{}) error) (Reader, error) {
	return f(loader)
}

var HeaderFactory = func(field string) (Reader, error) {
	return func(r *http.Request) ([]byte, error) {
		return []byte(r.Header.Get(field)), nil
	}, nil
}
var QueryFactory = func(field string) (Reader, error) {
	return func(r *http.Request) ([]byte, error) {
		q := r.URL.Query()
		return []byte(q.Get(field)), nil
	}, nil
}

var FormFactory = func(field string) (Reader, error) {
	return func(r *http.Request) ([]byte, error) {
		f := r.Form
		return []byte(f.Get(field)), nil
	}, nil
}
