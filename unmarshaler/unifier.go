package unmarshaler

import (
	"reflect"
	"strings"
)

type Unifier interface {
	Unify(a *Assembler, rv reflect.Value) (bool, error)
}

type Unifiers map[Type][]Unifier

func (u *Unifiers) Unify(a *Assembler, v interface{}) (bool, error) {
	rv := reflect.Indirect(reflect.ValueOf(v))

	return u.UnifyReflectValue(a, rv)
}
func (u *Unifiers) UnifyReflectValue(a *Assembler, rv reflect.Value) (bool, error) {
	// rv = reflect.Indirect(rv)
	tp, err := a.CheckType(rv.Type())
	if err != nil {
		return false, err
	}
	if tp == TypeUnkonwn {
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
func (u *Unifiers) Append(tp Type, unifier Unifier) {
	m := (*u)
	v := m[tp]
	v = append(v, unifier)
	m[tp] = v
	*u = m
}
func (u *Unifiers) Insert(tp Type, unifier Unifier) {
	m := (*u)
	v := []Unifier{unifier}
	v = append(v, m[tp]...)
	m[tp] = v
	*u = m

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
		rv.SetBool(s)
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
		rv.Set(reflect.ValueOf(s))
		return true, nil
	}
	if !a.Config().DisableConvertStringInterface {
		i, ok := v.(String)
		if ok {
			rv.Set(reflect.ValueOf(i))
			return true, nil
		}
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
			v := reflect.New(rv.Type().Elem()).Elem()
			_, err = a.Config().Unifiers.UnifyReflectValue(a.WithChild(iter.Part, rv.Type(), iter.Step), v)
			if err != nil {
				return false, err
			}
			sv = reflect.Append(sv, v)
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
		pv, err := iter.Part.Value()
		if err != nil {
			return false, err
		}
		mv.SetMapIndex(reflect.ValueOf(iter.Step.Interface()), reflect.ValueOf(pv))
		iter, err = iter.Next()
		if err != nil {
			return false, err
		}
	}
	rv.Set(mv)
	return true, nil
})

var UnifierEmptyInterface = UnifierFunc(func(a *Assembler, rv reflect.Value) (bool, error) {
	v, err := a.Part().Value()
	if err != nil {
		return false, err
	}
	rt := reflect.TypeOf(v)
	switch rt.Kind() {
	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8, reflect.String, reflect.Bool:
		rv.Set(reflect.ValueOf(v))
		return true, nil
	case reflect.Slice:
		return UnifierSlice(a, rv)
	case reflect.Map:
		return UnifierMap(a, rv)
	case reflect.Struct:
		return UnifierStruct(a, rv)
	}
	return false, nil
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
	value := reflect.New(rt).Elem()
	for i := 0; i < fl; i++ {
		var part Part
		var ok bool
		field := rt.Field(i)
		if field.PkgPath != "" {
			continue
		}
		fv := value.Field(i)
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

var UnifierLazyLoadFunc = UnifierFunc(func(a *Assembler, rv reflect.Value) (bool, error) {
	l := NewLazyLoader()
	l.Assembler = a
	rv.Set(reflect.ValueOf(l.LazyLoad))
	return true, nil
})

var UnifierLazyLoader = UnifierFunc(func(a *Assembler, rv reflect.Value) (bool, error) {
	l := NewLazyLoader()
	l.Assembler = a
	rv.Set(reflect.ValueOf(l))
	return true, nil
})

var UnifierPtr = UnifierFunc(func(a *Assembler, rv reflect.Value) (bool, error) {
	v := reflect.New(rv.Type().Elem())
	rv.Set(v)
	return a.Config().Unifiers.UnifyReflectValue(a, v.Elem())
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
	u.Append(TypeStruct, UnifierStruct)
	u.Append(TypePtr, UnifierPtr)
	u.Append(TypeEmptyInterface, UnifierEmptyInterface)
	u.Append(TypeLazyLoadFunc, UnifierLazyLoadFunc)
	u.Append(TypeLazyLoader, UnifierLazyLoader)

}
