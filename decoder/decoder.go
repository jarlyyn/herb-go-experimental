package decoder

import (
	"reflect"
)

type Decoder struct {
	Checkers []TypeChecker
	Unifiers Unifiers
}

func (d *Decoder) DecodeDataSource(v interface{}, n Node) error {
	ctx := NewContext()
	ctx.Decoder = d
	ctx.Node = n
	data, err := n.GetData(nil)
	if err != nil {
		return err
	}
	_, err = d.Unifiers.Unify(ctx, reflect.ValueOf(v), data)
	return err
}

func (d *Decoder) CheckType(ctx *Context, rt reflect.Type) (tp interface{}, err error) {
	for k := range d.Checkers {
		result, err := d.Checkers[k].CheckType(ctx, rt)
		if err != nil {
			return nil, err
		}
		if result {
			return d.Checkers[k].Type, nil
		}
	}
	return nil, nil
}

func NewDecoder() *Decoder {
	return &Decoder{}
}

func NewCommonDecoder() *Decoder {
	d := NewDecoder()
	return d
}
