package requestparamreader

import (
	"net"
	"net/http"

	"github.com/herb-go/herb/middleware/router"
)

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

var RouterFactory = newCommonFactory(func(r *http.Request, field string) ([]byte, error) {
	p := router.GetParams(r)
	return []byte(p.Get(field)), nil
})
var FixedFactory = newCommonFactory(func(r *http.Request, field string) ([]byte, error) {
	return []byte(field), nil
})
var CookieFactory = newCommonFactory(func(r *http.Request, field string) ([]byte, error) {
	c, err := r.Cookie(field)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, nil
		}
		return nil, err
	}
	return []byte(c.Value), nil
})
var IPAddressFactory = ReaderFactoryFunc(func(loader func(v interface{}) error) (Reader, error) {
	return func(r *http.Request) ([]byte, error) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return nil, err
		}
		return []byte(ip), nil
	}, nil
})

var MethodFactory = ReaderFactoryFunc(func(loader func(v interface{}) error) (Reader, error) {
	return func(r *http.Request) ([]byte, error) {
		return []byte(r.Method), nil
	}, nil
})

var PathFactory = ReaderFactoryFunc(func(loader func(v interface{}) error) (Reader, error) {
	return func(r *http.Request) ([]byte, error) {
		return []byte(r.URL.Path), nil
	}, nil
})

var HostFactory = ReaderFactoryFunc(func(loader func(v interface{}) error) (Reader, error) {
	return func(r *http.Request) ([]byte, error) {
		return []byte(r.Host), nil
	}, nil
})

var UserFactory = ReaderFactoryFunc(func(loader func(v interface{}) error) (Reader, error) {
	return func(r *http.Request) ([]byte, error) {
		u, _, _ := r.BasicAuth()
		return []byte(u), nil
	}, nil
})

var PasswordFactory = ReaderFactoryFunc(func(loader func(v interface{}) error) (Reader, error) {
	return func(r *http.Request) ([]byte, error) {
		_, p, _ := r.BasicAuth()
		return []byte(p), nil
	}, nil
})
