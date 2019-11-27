package unmarshaler

import (
	"reflect"
	"sync"
)

type Config struct {
	Checkers                      *TypeCheckers
	Unifiers                      *Unifiers
	TagName                       string
	TagLazyLoad                   string
	TagParser                     func(c *Config, value string) (*Tag, error)
	CaseSensitive                 bool
	DisableConvertStringInterface bool
	DisableConvertNumber          bool
	CachedTags                    sync.Map
}

func (c *Config) GetTags(structType reflect.Type, field reflect.StructField) (*Tag, error) {
	if c.TagName == "" {
		return nil, nil
	}
	sname := structType.Name()
	if sname != "" {
		key := structType.PkgPath() + "." + structType.Name() + "." + field.Name
		t, ok := c.CachedTags.Load(key)
		if ok {
			return t.(*Tag), nil
		}
		tag, err := c.TagParser(c, field.Tag.Get(c.TagName))
		if err != nil {
			return nil, err
		}
		c.CachedTags.Store(key, tag)
		return tag, nil
	}
	return c.TagParser(c, field.Tag.Get(c.TagName))
}

func NewConfig() *Config {
	return &Config{
		TagParser:   ParseTag,
		Unifiers:    &Unifiers{},
		Checkers:    &TypeCheckers{},
		TagLazyLoad: "lazyload",
	}
}

func NewCommonConfig() *Config {
	c := NewConfig()
	c.TagName = "config"
	SetCommonTypeCheckers(c.Checkers)
	SetCommonUnifiers(c.Unifiers)
	return c
}
