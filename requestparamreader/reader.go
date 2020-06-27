package requestparamreader

import "net/http"

type Reader func(r *http.Request) ([]byte, error)

type ReaderFactory interface {
	CreateReader(loader func(v interface{}) error) (Reader, error)
}

type ReaderFactoryFunc func(loader func(v interface{}) error) (Reader, error)

func (f ReaderFactoryFunc) CreateReader(loader func(v interface{}) error) (Reader, error) {
	return f(loader)
}

func newCommonFactory(fieldloader func(r *http.Request, field string) ([]byte, error)) ReaderFactory {
	return ReaderFactoryFunc(func(loader func(v interface{}) error) (Reader, error) {
		c := &CommonFieldConfig{}
		err := loader(c)
		if err != nil {
			return nil, err
		}
		return func(r *http.Request) ([]byte, error) {
			field, err := fieldloader(r, c.Field)
			if err != nil {
				return nil, err
			}
			return field, nil
		}, nil

	})
}

var HeaderFactory = newCommonFactory(func(r *http.Request, field string) ([]byte, error) {
	return []byte(r.Header.Get(field)), nil
})
var QueryFactory = newCommonFactory(func(r *http.Request, field string) ([]byte, error) {
	q := r.URL.Query()
	return []byte(q.Get(field)), nil
})

var FormFactory = newCommonFactory(func(r *http.Request, field string) ([]byte, error) {
	f := r.Form
	return []byte(f.Get(field)), nil
})
