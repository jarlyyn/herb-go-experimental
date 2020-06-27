package requestparamreader

import (
	"errors"
	"sync"
)

var lock sync.Mutex
var builtinReaderFactories = map[string]ReaderFactory{}
var registeredReaderFactories = map[string]ReaderFactory{}
var builtinConverterFactories = map[string]ConverterFactory{}
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

func initBuiltinReaderFactories() {
	builtinReaderFactories["header"] = HeaderReaderFactory
	builtinReaderFactories["query"] = QueryReaderFactory
	builtinReaderFactories["form"] = FormReaderFactory
	builtinReaderFactories["router"] = RouterReaderFactory
	builtinReaderFactories["fixed"] = FixedReaderFactory
	builtinReaderFactories["cookie"] = CookieReaderFactory
	builtinReaderFactories["ip"] = IPAddressReaderFactory
	builtinReaderFactories["method"] = MethodReaderFactory
	builtinReaderFactories["path"] = PathReaderFactory
	builtinReaderFactories["host"] = HostReaderFactory
	builtinReaderFactories["user"] = UserReaderFactory
	builtinReaderFactories["passwrod"] = PasswordReaderFactory

}

func GetReaderFactoryByName(name string) (ReaderFactory, error) {
	f, ok := registeredReaderFactories[name]
	if ok {
		return f, nil
	}
	f, ok = builtinReaderFactories[name]
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
	f, ok = builtinConverterFactories[name]
	if ok {
		return f, nil
	}
	return nil, ErrFactoryNotRegistered
}
