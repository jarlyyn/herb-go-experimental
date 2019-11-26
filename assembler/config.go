package assembler

import (
	"reflect"
	"sync"
)

type Config struct {
	Checkers      *TypeCheckers
	Unifiers      *Unifiers
	TagName       string
	TagParser     func(value string) (*Tag, error)
	CaseSensitive bool
	CachedTags    sync.Map
}

func (c *Config) GetTags(structType reflect.Type, field reflect.StructField) (*Tag, error) {
	if c.TagName == "" {
		return nil, nil
	}
	key := structType.PkgPath() + "." + structType.Name() + "." + field.Name
	t, ok := c.CachedTags.Load(key)
	if ok {
		return t.(*Tag), nil
	}

	return c.TagParser(field.Tag.Get(c.TagName))
}

func NewConfig() *Config {
	return &Config{
		TagParser: ParseTag,
		Unifiers:  &Unifiers{},
		Checkers:  &TypeCheckers{},
	}
}

func NewCommonConfig() *Config {
	c := NewConfig()
	c.TagName = "config"
	SetCommonTypeCheckers(c.Checkers)
	SetCommonUnifiers(c.Unifiers)
	return c
}
