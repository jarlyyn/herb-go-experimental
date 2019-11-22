package decoder

import (
	"reflect"
)

type Context struct {
	Decoder    *Decoder
	DataSource DataSource
	Path       Path
	Parent     reflect.Type
}

func (c *Context) Clone() *Context {
	return &Context{
		Decoder:    c.Decoder,
		DataSource: c.DataSource,
		Path:       c.Path.Clone(),
		Parent:     c.Parent,
	}
}

func (c *Context) CloneAndJoin(step ...Step) (*Context, interface{}, error) {
	ctx := c.Clone()
	ctx.Path.Join(step...)
	v, err := ctx.DataSource.GetData(ctx.Path)
	if err != nil {
		return nil, nil, err
	}
	return ctx, v, nil
}

func NewContext() *Context {
	return &Context{
		Path: NewSteps(),
	}
}
