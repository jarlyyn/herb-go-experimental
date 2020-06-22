package requestparamreader

import (
	"errors"
	"sync"
)

var lock sync.Mutex
var buildinReaderFactories = map[string]ReaderFactory{}
var registeredReaderFactories = map[string]ReaderFactory{}
var buildinConverterFactories = map[string]ConverterFactory{}
var registeredConverterFactories = map[string]ConverterFactory{}

var ErrFactoryNotRegistered = errors.New("err factory not registered")
var ErrUnavaliableField = errors.New("unavaliable field")

func Reset() {
	lock.Lock()
	defer lock.Unlock()
	registeredReaderFactories = map[string]ReaderFactory{}
	registeredConverterFactories = map[string]ConverterFactory{}
}

func RegisterReaderFactory(name string, f ReaderFactory) {
	lock.Lock()
	defer lock.Unlock()
	registeredReaderFactories[name] = f
}

func GetReaderFactoryByName(name string) (ReaderFactory, error) {
	f, ok := registeredReaderFactories[name]
	if ok {
		return f, nil
	}
	f, ok = buildinReaderFactories[name]
	if ok {
		return f, nil
	}
	return nil, ErrFactoryNotRegistered
}

func RegisterConverterFactory(name string, f ConverterFactory) {
	lock.Lock()
	defer lock.Unlock()
	registeredConverterFactories[name] = f
}

func GetConverterFactoryByName(name string) (ConverterFactory, error) {
	f, ok := registeredConverterFactories[name]
	if ok {
		return f, nil
	}
	f, ok = buildinConverterFactories[name]
	if ok {
		return f, nil
	}
	return nil, ErrFactoryNotRegistered
}
