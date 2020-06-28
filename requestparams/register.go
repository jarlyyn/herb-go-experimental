package requestparam

import (
	"errors"
	"fmt"
	"sync"
)

var lock sync.Mutex
var builtinReaderFactories = map[string]ReaderFactory{}
var builtinFormaterFactories = map[string]FormaterFactory{}
var registeredReaders = map[string]Reader{}
var registeredFormaters = map[string]Formater{}
var ErrFactoryNotRegistered = errors.New("err factory not registered")
var ErrUnavaliableField = errors.New("unavaliable field")

func Reset() {
	lock.Lock()
	defer lock.Unlock()
	registeredReaders = map[string]Reader{}
	registeredFormaters = map[string]Formater{}
}

func RegisterReader(name string, reader Reader) {
	lock.Lock()
	defer lock.Unlock()
	registeredReaders[name] = reader
}
func GetReader(name string) (Reader, error) {
	lock.Lock()
	defer lock.Unlock()
	r := registeredReaders[name]
	if r == nil {
		return nil, fmt.Errorf("%W :%S", ErrReaderNotFound, name)
	}
	return r, nil
}

func GetFormater(name string) (Formater, error) {
	lock.Lock()
	defer lock.Unlock()
	f := registeredFormaters[name]
	if f == nil {
		return nil, fmt.Errorf("%W :%S", ErrFormaterNotFound, name)
	}
	return f, nil
}

func RegisterFormater(name string, formater Formater) {
	lock.Lock()
	defer lock.Unlock()
	registeredFormaters[name] = formater
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
func initBuiltinFormatFactories() {
	builtinFormaterFactories["toupper"] = ToUpperFormaterFactory
	builtinFormaterFactories["tolower"] = ToLowerFormaterFactory
	builtinFormaterFactories["trim"] = TrimSpaceFormaterFactory
	builtinFormaterFactories["interger"] = IntegerFormaterFactory
	builtinFormaterFactories["match"] = MatchFormaterFactory
	builtinFormaterFactories["find"] = FindFormaterFactory

}
func GetReaderFactoryByName(name string) (ReaderFactory, error) {

	f, ok := builtinReaderFactories[name]
	if ok {
		return f, nil
	}
	return nil, ErrFactoryNotRegistered
}

func GetFormaterFactoryByName(name string) (FormaterFactory, error) {
	f, ok := builtinFormaterFactories[name]
	if ok {
		return f, nil
	}
	return nil, ErrFactoryNotRegistered
}

func init() {
	initBuiltinReaderFactories()
	initBuiltinFormatFactories()
}
