package assembler

import (
	"reflect"
	"strings"
)

type Unifier interface {
	Unify(a *Assembler, rv reflect.Value) (bool, error)
}

type Unifiers map[interface{}][]Unifier

func (u *Unifiers) Unify(a *Assembler, v interface{}) (bool, error) {
	rv := reflect.Indirect(reflect.ValueOf(v))
	return u.UnifyReflectValue(a, rv)
}
func (u *Unifiers) UnifyReflectValue(a *Assembler, rv reflect.Value) (bool, error) {
	tp, err := a.CheckType()
	if err != nil {
		return false, err
	}
	if tp == nil {
		return false, nil
	}
	unifiers, ok := (*u)[tp]
	if ok == false {
		return false, nil
	}
	for k := range unifiers {
		result, err := unifiers[k].Unify(a, rv)
		if err != nil {
			return false, err
		}
		if result {
			return true, nil
		}
	}
	return false, nil

}
func (u *Unifiers) Append(tp interface{}, unifier Unifier) {
	v := (*u)[tp]
	v = append(v, unifier)
	(*u)[tp] = v
}
func (u *Unifiers) Insert(tp interface{}, unifier Unifier) {
	v := []Unifier{unifier}
	v = append(v, (*u)[tp]...)
	(*u)[tp] = v
}

type String interface {
	String() string
}
type UnifierFunc func(a *Assembler, rv reflect.Value) (bool, error)

func (f UnifierFunc) Unify(a *Assembler, rv reflect.Value) (bool, error) {
	return f(a, rv)
}

var UnifierBool = UnifierFunc(func(a *Assembler, rv reflect.Value) (bool, error) {
	v, err := a.Part().Value()
	if err != nil {
		return false, err
	}
	s, ok := v.(bool)
	if ok {
		reflect.ValueOf(v).SetBool(s)
		return true, nil
	}
	i, ok := v.(String)
	if ok {
		rv.Set(reflect.ValueOf(i))
		return true, nil
	}
	return false, nil
})

var UnifierString = UnifierFunc(func(a *Assembler, rv reflect.Value) (bool, error) {
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
		rv.Set(reflect.ValueOf(i))
		return true, nil
	}
	return false, nil
})

var UnifierInt = UnifierFunc(func(a *Assembler, rv reflect.Value) (bool, error) {
	v, err := a.Part().Value()
	if err != nil {
		return false, err
	}
	s, ok := v.(int)
	if ok {
		rv.Set(reflect.ValueOf(s))
		return true, nil
	}
	return false, nil
})

var UnifierUint = UnifierFunc(func(a *Assembler, rv reflect.Value) (bool, error) {
	v, err := a.Part().Value()
	if err != nil {
		return false, err
	}
	s, ok := v.(uint)
	if ok {
		rv.Set(reflect.ValueOf(s))
		return true, nil
	}
	return false, nil
})

var UnifierInt64 = UnifierFunc(func(a *Assembler, rv reflect.Value) (bool, error) {
	v, err := a.Part().Value()
	if err != nil {
		return false, err
	}
	s, ok := v.(int64)
	if ok {
		rv.Set(reflect.ValueOf(s))
		return true, nil
	}
	return false, nil
})
var UnifierUint64 = UnifierFunc(func(a *Assembler, rv reflect.Value) (bool, error) {
	v, err := a.Part().Value()
	if err != nil {
		return false, err
	}
	s, ok := v.(uint64)
	if ok {
		rv.Set(reflect.ValueOf(s))
		return true, nil
	}
	return false, nil
})

var UnifierFloat32 = UnifierFunc(func(a *Assembler, rv reflect.Value) (bool, error) {
	v, err := a.Part().Value()
	if err != nil {
		return false, err
	}
	s, ok := v.(float32)
	if ok {
		rv.Set(reflect.ValueOf(s))
		return true, nil
	}
	return false, nil
})

var UnifierFloat64 = UnifierFunc(func(a *Assembler, rv reflect.Value) (bool, error) {
	v, err := a.Part().Value()
	if err != nil {
		return false, err
	}
	s, ok := v.(float64)
	if ok {
		rv.Set(reflect.ValueOf(s))
		return true, nil
	}
	return false, nil
})

var UnifierSlice = UnifierFunc(func(a *Assembler, rv reflect.Value) (bool, error) {
	iter, err := a.Part().Iter()
	if err != nil {
		return false, err
	}
	if iter == nil {
		return false, nil
	}
	sv := reflect.MakeSlice(rv.Type(), 0, 0)
	for iter != nil {
		if iter.Step.Type() == TypeArray {
			pv, err := iter.Part.Value()
			if err != nil {
				return false, err
			}
			rv = reflect.Append(rv, reflect.ValueOf(pv))
		}
		iter, err = iter.Next()
		if err != nil {
			return false, err
		}
	}
	rv.Set(sv)
	return true, nil
})

var UnifierMap = UnifierFunc(func(a *Assembler, rv reflect.Value) (bool, error) {
	iter, err := a.Part().Iter()
	if err != nil {
		return false, err
	}
	if iter == nil {
		return false, nil
	}

	mv := reflect.MakeMap(rv.Type())
	for iter != nil {
		if iter.Step.Type() == TypeArray {
			pv, err := iter.Part.Value()
			if err != nil {
				return false, err
			}
			rv.SetMapIndex(reflect.ValueOf(iter.Step.Interface()), reflect.ValueOf(pv))
		}
		iter, err = iter.Next()
		if err != nil {
			return false, err
		}
	}
	rv.Set(mv)
	return true, nil
})
var UnifierStruct = UnifierFunc(func(a *Assembler, rv reflect.Value) (bool, error) {
	iter, err := a.Part().Iter()
	if err != nil {
		return false, err
	}
	if iter == nil {
		return false, nil
	}
	var valuemap = map[string]Part{}
	var civaluemap = map[string]Part{}
	var part Part
	var ok bool
	ci := !a.Config().CaseSensitive
	for iter != nil {

		valuemap[iter.Step.String()] = iter.Part
		if ci {
			civaluemap[strings.ToLower(iter.Step.String())] = iter.Part
		}
		iter, err = iter.Next()
		if err != nil {
			return false, err
		}
	}
	rt := rv.Type()
	fl := rt.NumField()
	value := reflect.New(rt)
	for i := 0; i < fl; i++ {
		field := rt.FieldByIndex([]int{i})
		fv := value.FieldByIndex([]int{i})
		tag, err := a.Config().GetTags(rt, field)
		if err != nil {
			return false, err
		}
		if tag.Name != "" {
			part, ok = valuemap[tag.Name]
		}
		if !ok {
			part, ok = valuemap[field.Name]
		}
		if !ok && ci {
			part, ok = civaluemap[strings.ToLower(field.Name)]
		}
		if !ok {
			continue
		}
		_, err = a.Config().Unifiers.UnifyReflectValue(a.WithChild(part, rt, NewFieldStep(&field)), fv)
		if err != nil {
			return false, err
		}
	}
	rv.Set(value)
	return true, nil
})

func SetCommonUnifiers(u *Unifiers) {
	u.Append(TypeBool, UnifierBool)
	u.Append(TypeString, UnifierString)
	u.Append(TypeInt, UnifierInt)
	u.Append(TypeUint, UnifierUint)
	u.Append(TypeInt64, UnifierInt64)
	u.Append(TypeUint64, UnifierUint64)
	u.Append(TypeFloat32, UnifierFloat32)
	u.Append(TypeFloat64, UnifierFloat64)
	u.Append(TypeSlice, UnifierSlice)
	u.Append(TypeMap, UnifierMap)
	u.Append(TypeMap, UnifierStruct)
}
