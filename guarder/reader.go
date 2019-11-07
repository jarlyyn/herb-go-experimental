package guarder

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"sync"
)

type RequestParamsReader interface {
	ReadParamsFromRequest(*http.Request) (*RequestParams, error)
}

type RequestParamsReaderDriverField struct {
	ReaderDriver       string
	staticReaderDriver string
}

func (f *RequestParamsReaderDriverField) RequestParamsReaderDriver() string {
	if f.staticReaderDriver == "" {
		return f.ReaderDriver
	}
	return f.staticReaderDriver
}

//ReaderFactory guarder factory
type ReaderFactory func(conf Config, prefix string) (RequestParamsReader, error)

var (
	readerFactorysMu sync.RWMutex
	readerFactories  = make(map[string]ReaderFactory)
)

// RegisterReader makes a driver creator available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func RegisterReader(name string, f ReaderFactory) {
	readerFactorysMu.Lock()
	defer readerFactorysMu.Unlock()
	if f == nil {
		panic(errors.New("guarder: Register reader factory is nil"))
	}
	if _, dup := readerFactories[name]; dup {
		panic(errors.New("guarder: Register called twice for reader factory " + name))
	}
	readerFactories[name] = f
}

//UnregisterAllReader unregister all driver
func UnregisterAllReader() {
	readerFactorysMu.Lock()
	defer readerFactorysMu.Unlock()
	// For tests.
	readerFactories = make(map[string]ReaderFactory)
}

//ReaderFactories returns a sorted list of the names of the registered factories.
func ReaderFactories() []string {
	readerFactorysMu.RLock()
	defer readerFactorysMu.RUnlock()
	var list []string
	for name := range readerFactories {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

//NewReaderDriver create new driver with given name,config and prefix.
//Reutrn driver created and any error if raised.
func NewReaderDriver(name string, conf Config, prefix string) (RequestParamsReader, error) {
	readerFactorysMu.RLock()
	factoryi, ok := readerFactories[name]
	readerFactorysMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("guarder: unknown reader driver %q (forgotten import?)", name)
	}
	return factoryi(conf, prefix)
}
