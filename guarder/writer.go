package guarder

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"sync"
)

type RequestParamsWriter interface {
	WriterParamsToRequest(*http.Request, *RequestParams) error
}

//WriterFactory guarder factory
type WriterFactory func(conf Config, prefix string) (RequestParamsWriter, error)

var (
	writerFactorysMu sync.RWMutex
	writerFactories  = make(map[string]WriterFactory)
)

// RegisterWriter makes a driver creator available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func RegisterWriter(name string, f WriterFactory) {
	writerFactorysMu.Lock()
	defer writerFactorysMu.Unlock()
	if f == nil {
		panic(errors.New("guarder: Register writer factory is nil"))
	}
	if _, dup := writerFactories[name]; dup {
		panic(errors.New("guarder: Register called twice for writer factory " + name))
	}
	writerFactories[name] = f
}

//UnregisterAllWriter unregister all driver
func UnregisterAllWriter() {
	writerFactorysMu.Lock()
	defer writerFactorysMu.Unlock()
	// For tests.
	writerFactories = make(map[string]WriterFactory)
}

//WriterFactories returns a sorted list of the names of the registered factories.
func WriterFactories() []string {
	writerFactorysMu.RLock()
	defer writerFactorysMu.RUnlock()
	var list []string
	for name := range writerFactories {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

//NewWriterDriver create new driver with given name,config and prefix.
//Reutrn driver created and any error if raised.
func NewWriterDriver(name string, conf Config, prefix string) (RequestParamsWriter, error) {
	writerFactorysMu.RLock()
	factoryi, ok := writerFactories[name]
	writerFactorysMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("guarder: unknown writer driver %q (forgotten import?)", name)
	}
	return factoryi(conf, prefix)
}
