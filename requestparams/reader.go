package requestparam

import (
	"fmt"
	"net"
	"net/http"

	"github.com/herb-go/worker"

	"github.com/herb-go/providers/herb/overseers/identifieroverseer"

	"github.com/herb-go/herb/user/httpuser"

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

func newCommonReaderFactory(fieldloader func(r *http.Request, field string) ([]byte, error)) ReaderFactory {
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

var HeaderReaderFactory = newCommonReaderFactory(func(r *http.Request, field string) ([]byte, error) {
	return []byte(r.Header.Get(field)), nil
})
var QueryReaderFactory = newCommonReaderFactory(func(r *http.Request, field string) ([]byte, error) {
	q := r.URL.Query()
	return []byte(q.Get(field)), nil
})

var FormReaderFactory = newCommonReaderFactory(func(r *http.Request, field string) ([]byte, error) {
	f := r.Form
	return []byte(f.Get(field)), nil
})

var RouterReaderFactory = newCommonReaderFactory(func(r *http.Request, field string) ([]byte, error) {
	p := router.GetParams(r)
	return []byte(p.Get(field)), nil
})
var FixedReaderFactory = newCommonReaderFactory(func(r *http.Request, field string) ([]byte, error) {
	return []byte(field), nil
})
var CookieReaderFactory = newCommonReaderFactory(func(r *http.Request, field string) ([]byte, error) {
	c, err := r.Cookie(field)
	if err != nil {
		if err == http.ErrNoCookie {
			return nil, nil
		}
		return nil, err
	}
	return []byte(c.Value), nil
})
var IPAddressReaderFactory = ReaderFactoryFunc(func(loader func(v interface{}) error) (Reader, error) {
	return func(r *http.Request) ([]byte, error) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			return nil, err
		}
		return []byte(ip), nil
	}, nil
})

var MethodReaderFactory = ReaderFactoryFunc(func(loader func(v interface{}) error) (Reader, error) {
	return func(r *http.Request) ([]byte, error) {
		return []byte(r.Method), nil
	}, nil
})

var PathReaderFactory = ReaderFactoryFunc(func(loader func(v interface{}) error) (Reader, error) {
	return func(r *http.Request) ([]byte, error) {
		return []byte(r.URL.Path), nil
	}, nil
})

var HostReaderFactory = ReaderFactoryFunc(func(loader func(v interface{}) error) (Reader, error) {
	return func(r *http.Request) ([]byte, error) {
		return []byte(r.Host), nil
	}, nil
})

var UserReaderFactory = ReaderFactoryFunc(func(loader func(v interface{}) error) (Reader, error) {
	return func(r *http.Request) ([]byte, error) {
		u, _, _ := r.BasicAuth()
		return []byte(u), nil
	}, nil
})

var PasswordReaderFactory = ReaderFactoryFunc(func(loader func(v interface{}) error) (Reader, error) {
	return func(r *http.Request) ([]byte, error) {
		_, p, _ := r.BasicAuth()
		return []byte(p), nil
	}, nil
})

type IdentifierReader struct {
	Identifier httpuser.Identifier
}

func (reader *IdentifierReader) Read(r *http.Request) ([]byte, error) {
	id, err := reader.Identifier.IdentifyRequest(r)
	if err != nil {
		return nil, err
	}
	return []byte(id), nil
}

type WorkerConfig struct {
	ID string
}

var IdentifierReaderFactory = ReaderFactoryFunc(func(loader func(v interface{}) error) (Reader, error) {
	c := &WorkerConfig{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	identifier := identifieroverseer.GetIdentifierByID(c.ID)
	if identifier == nil {
		return nil, fmt.Errorf("%w :%s", worker.ErrWorkerNotFound, c.ID)
	}
	reader := &IdentifierReader{Identifier: identifier}
	return reader.Read, nil
})
