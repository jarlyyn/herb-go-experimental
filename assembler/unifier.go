package assembler

import (
	"reflect"
)

type Unifier interface {
	Unify(a *Assembler, v interface{}) (bool, error)
}

type Unifiers map[interface{}][]Unifier

func (u *Unifiers) Unify(a *Assembler, v interface{}) error {
	v = reflect.Indirect(reflect.ValueOf(v))
	tp, err := a.CheckType()
	if err != nil {
		return err
	}
	unifiers, ok := (*u)[tp]
	if ok == false {
		return nil
	}
	for k := range unifiers {
		result, err := unifiers[k].Unify(a, v)
		if err != nil {
			return err
		}
		if result {
			return nil
		}
	}
	return nil
}

type String interface {
	String() string
}
type UnifierFunc func(a *Assembler, v interface{}) (bool, error)

func (f UnifierFunc) Unify(a *Assembler, v interface{}) (bool, error) {
	return f(a, v)
}

var StringUnifier = UnifierFunc(func(a *Assembler, value interface{}) (bool, error) {
	v, err := a.Part().Value()
	if err != nil {
		return false, err
	}
	s, ok := v.(string)
	if ok {
		reflect.ValueOf(v).SetString(s)
		return true, nil
	}
	i, ok := v.(String)
	if ok {
		reflect.ValueOf(value).Set(reflect.ValueOf(i))
		return true, nil
	}
	return false, nil
})

var IntUnifier = UnifierFunc(func(a *Assembler, value interface{}) (bool, error) {
	v, err := a.Part().Value()
	if err != nil {
		return false, err
	}
	s, ok := v.(int)
	if ok {
		reflect.ValueOf(v).Set(reflect.ValueOf(s))
		return true, nil
	}
	return false, nil
})

var UintUnifier = UnifierFunc(func(a *Assembler, value interface{}) (bool, error) {
	v, err := a.Part().Value()
	if err != nil {
		return false, err
	}
	s, ok := v.(uint)
	if ok {
		reflect.ValueOf(v).Set(reflect.ValueOf(s))
		return true, nil
	}
	return false, nil
})

var Int64Unifier = UnifierFunc(func(a *Assembler, value interface{}) (bool, error) {
	v, err := a.Part().Value()
	if err != nil {
		return false, err
	}
	s, ok := v.(int64)
	if ok {
		reflect.ValueOf(v).Set(reflect.ValueOf(s))
		return true, nil
	}
	return false, nil
})
var Uint64Unifier = UnifierFunc(func(a *Assembler, value interface{}) (bool, error) {
	v, err := a.Part().Value()
	if err != nil {
		return false, err
	}
	s, ok := v.(uint64)
	if ok {
		reflect.ValueOf(v).Set(reflect.ValueOf(s))
		return true, nil
	}
	return false, nil
})
