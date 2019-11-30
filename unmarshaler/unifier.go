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

var UnifierNumber = UnifierFunc(func(a *Assembler, rv reflect.Value) (bool, error) {
	v, err := a.Part().Value()
	if err != nil {
		return false, err
	}
	av := reflect.ValueOf(v)
	switch av.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		if rv.Kind() == av.Kind() {
			rv.Set(av)
			return true, nil
		}
		rv.Set(reflect.ValueOf(v).Convert(rv.Type()))
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

func convertIterToArray(iter *PartIter) ([]interface{}, error) {
	a := []interface{}{}
	for iter != nil {
		pv, err := iter.Part.Value()
		if err != nil {
			return nil, err
		}
		a = append(a, pv)
		iter, err = iter.Next()
		if err != nil {
			return nil, err
		}
	}
	return a, nil
}

func convertIterToStringMap(iter *PartIter) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	for iter != nil {
		pv, err := iter.Part.Value()
		if err != nil {
			return nil, err
		}
		m[iter.Step.String()] = pv
		iter, err = iter.Next()
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}
func convertIterToInterfaceMap(iter *PartIter) (map[interface{}]interface{}, error) {
	m := map[interface{}]interface{}{}
	for iter != nil {
		pv, err := iter.Part.Value()
		if err != nil {
			return nil, err
		}
		m[iter.Step.Interface()] = pv
		iter, err = iter.Next()
		if err != nil {
			return nil, err
		}
	}
	return m, nil
}
func convertIter(i *PartIter) (interface{}, error) {
	switch i.Step.Type() {
	case TypeArray:
		return convertIterToArray(i)
	case TypeString:
		return convertIterToStringMap(i)
	case TypeEmptyInterface:
		return convertIterToInterfaceMap(i)
	}
	return nil, nil
}

var UnifierEmptyInterface = UnifierFunc(func(a *Assembler, rv reflect.Value) (bool, error) {
	iter, err := a.Part().Iter()
	if err != nil {
		return false, err
	}
	if iter == nil {
		v, err := a.Part().Value()
		if err != nil {
			return false, err
		}
		rt := reflect.TypeOf(v)
		switch rt.Kind() {
		case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
			reflect.String, reflect.Bool,
			reflect.Map, reflect.Slice:
			rv.Set(reflect.ValueOf(v))
			return true, nil
		}
	} else {
		val, err := convertIter(iter)
		if err != nil {
			return false, err
		}
		if val == nil {
			return false, nil
		}
		rv.Set(reflect.ValueOf(val))
	}
	return false, nil
})

type structData struct {
	assembler  *Assembler
	valuemap   map[string]Part
	civaluemap map[string]Part
}

func (d *structData) LoadValues() (bool, error) {
	a := d.assembler
	iter, err := a.Part().Iter()
	if err != nil {
		return false, err
	}
	if iter == nil {
		return false, nil
	}
	d.valuemap = map[string]Part{}
	d.civaluemap = map[string]Part{}
	ci := !a.Config().CaseSensitive
	for iter != nil {

		d.valuemap[iter.Step.String()] = iter.Part
		if ci {
			d.civaluemap[strings.ToLower(iter.Step.String())] = iter.Part
		}
		iter, err = iter.Next()
		if err != nil {
			return false, err
		}
	}
	return true, nil
}
func (d *structData) IsAnonymous(field reflect.StructField, tag *Tag) bool {
	if field.Type.Kind() != reflect.Struct {
		return false
	}
	if tag.Ignored || tag.Name != "" {
		return false
	}
	c := d.assembler.Config()
	if c.TagAnonymous != "" && tag.Flags[c.TagAnonymous] != "" {
		return true
	}
	if d.valuemap[field.Name] != nil {
		return false
	}
	ci := !c.CaseSensitive
	if ci && d.civaluemap[strings.ToLower(field.Name)] != nil {
		return false
	}
	return true
}
func (d *structData) WalkStruct(rv reflect.Value) (bool, error) {
	a := d.assembler
	rt := rv.Type()
	fl := rt.NumField()
	ci := !a.Config().CaseSensitive
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
		if tag.Ignored {
			continue
		}
		if d.IsAnonymous(field, tag) {
			_, err := d.WalkStruct(fv)
			if err != nil {
				return false, err
			}
			continue
		}
		if err != nil {
			return false, err
		}
		if tag.Name != "" {
			part, ok = d.valuemap[tag.Name]
		}
		if !ok {
			part, ok = d.valuemap[field.Name]
		}
		if !ok && ci {
			part, ok = d.civaluemap[strings.ToLower(field.Name)]
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

}
func newStructData() *structData {
	return &structData{
		valuemap:   map[string]Part{},
		civaluemap: map[string]Part{},
	}
}

var UnifierStruct = UnifierFunc(func(a *Assembler, rv reflect.Value) (bool, error) {
	sd := newStructData()
	sd.assembler = a
	ok, err := sd.LoadValues()
	if ok == false || err != nil {
		return ok, err
	}
	return sd.WalkStruct(rv)
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
	u.Append(TypeInt, UnifierNumber)
	u.Append(TypeUint, UnifierNumber)
	u.Append(TypeInt64, UnifierNumber)
	u.Append(TypeUint64, UnifierNumber)
	u.Append(TypeFloat32, UnifierNumber)
	u.Append(TypeFloat64, UnifierNumber)
	u.Append(TypeSlice, UnifierSlice)
	u.Append(TypeMap, UnifierMap)
	u.Append(TypeStruct, UnifierStruct)
	u.Append(TypePtr, UnifierPtr)
	u.Append(TypeEmptyInterface, UnifierEmptyInterface)
	u.Append(TypeLazyLoadFunc, UnifierLazyLoadFunc)
	u.Append(TypeLazyLoader, UnifierLazyLoader)

}
