package requestparam

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type Formater func([]byte) ([]byte, bool, error)

type FormaterFactory interface {
	CreateFormater(loader func(interface{}) error) (Formater, error)
}

type FormaterFactoryFunc func(loader func(interface{}) error) (Formater, error)

func (f FormaterFactoryFunc) CreateFormater(loader func(interface{}) error) (Formater, error) {
	return f(loader)
}

var ToUpperFormaterFactory = FormaterFactoryFunc(func(loader func(interface{}) error) (Formater, error) {
	return func(data []byte) ([]byte, bool, error) {
		return []byte(strings.ToUpper(string(data))), true, nil
	}, nil
})

var ToLowerFormaterFactory = FormaterFactoryFunc(func(loader func(interface{}) error) (Formater, error) {
	return func(data []byte) ([]byte, bool, error) {
		return []byte(strings.ToLower(string(data))), true, nil
	}, nil
})

var TrimSpaceFormaterFactory = FormaterFactoryFunc(func(loader func(interface{}) error) (Formater, error) {
	return func(data []byte) ([]byte, bool, error) {
		return []byte(strings.TrimSpace(string(data))), true, nil
	}, nil
})

var IntegerFormaterFactory = FormaterFactoryFunc(func(loader func(interface{}) error) (Formater, error) {
	return func(data []byte) ([]byte, bool, error) {
		_, err := strconv.Atoi(string(data))
		if err != nil {
			return nil, false, nil
		}
		return []byte(strings.TrimSpace(string(data))), true, nil
	}, nil
})

type RegexpFormater struct {
	*regexp.Regexp
	Index int
}

type RegexpConfig struct {
	Pattern string
	Index   int
}

func (f *RegexpFormater) Match(data []byte) ([]byte, bool, error) {
	ok := f.Regexp.Match(data)
	if ok {
		return data, true, nil
	}
	return nil, false, nil
}

func (f *RegexpFormater) Find(data []byte) ([]byte, bool, error) {
	results := f.Regexp.FindSubmatch(data)
	if len(results) > f.Index {
		return results[f.Index], true, nil
	}
	return nil, false, nil
}

func getRegexp(loader func(interface{}) error) (*RegexpFormater, error) {
	c := &RegexpConfig{}
	err := loader(c)
	if err != nil {
		return nil, err
	}
	if c.Index < 0 {
		return nil, errors.New("unavailable submatch index")
	}
	p, err := regexp.Compile(c.Pattern)
	if err != nil {
		return nil, err
	}

	f := &RegexpFormater{
		Regexp: p,
		Index:  c.Index,
	}
	return f, nil

}

var MatchFormaterFactory = FormaterFactoryFunc(func(loader func(interface{}) error) (Formater, error) {
	f, err := getRegexp(loader)
	if err != nil {
		return nil, err
	}
	return f.Match, nil
})

var FindFormaterFactory = FormaterFactoryFunc(func(loader func(interface{}) error) (Formater, error) {
	f, err := getRegexp(loader)
	if err != nil {
		return nil, err
	}
	return f.Find, nil
})
