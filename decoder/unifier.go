package decoder

import "reflect"

type Unifier interface {
	Unify(ctx *Context, rv reflect.Value, v interface{}) (bool, error)
}

type Unifiers map[interface{}][]Unifier

func (u *Unifiers) Unify(ctx *Context, rv reflect.Value, v interface{}) (bool, error) {
	rv = reflect.Indirect(rv)
	tp, err := ctx.Decoder.CheckType(ctx, rv.Type())
	if err != nil {
		return false, err
	}
	unifiers, ok := (*u)[tp]
	if ok == false {
		return false, nil
	}
	for k := range unifiers {
		result, err := unifiers[k].Unify(ctx, rv, v)
		if err != nil {
			return false, err
		}
		if result {
			return result, nil
		}
	}
	return false, nil
}
