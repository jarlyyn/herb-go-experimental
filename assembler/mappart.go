package assembler

import (
	"reflect"
)

type MapPart struct {
	value interface{}
}

func NewMapPart(v interface{}) *MapPart {
	return &MapPart{
		value: v,
	}
}
func (d *MapPart) Value() (interface{}, error) {
	return d.value, nil
}
func (d *MapPart) mapIter(rt reflect.Type, rv reflect.Value) (*PartIter, error) {
	mes := mapElements{}
	iter := rv.MapRange()
	keykind := rt.Key().Kind()
	for iter.Next() {
		var key Step
		switch keykind {
		case reflect.String:
			key = NewStringStep(iter.Key().String())
		default:
			key = NewInterfaceStep(iter.Key().Interface())
		}
		m := mapElement{
			Step: key,
			Part: NewMapPart(iter.Value().Interface()),
		}

		mes = append(mes, m)
	}
	return mes.Next()
}
func (d *MapPart) arrayIter(rv reflect.Value) (*PartIter, error) {
	mes := mapElements{}
	l := rv.Len()
	for i := 0; i < l; i++ {
		m := mapElement{
			Step: NewArrayStep(i),
			Part: NewMapPart(rv.Index(i).Interface()),
		}
		mes = append(mes, m)
	}
	return mes.Next()
}
func (d *MapPart) Iter() (*PartIter, error) {
	v, err := d.Value()
	if err != nil {
		return nil, err
	}
	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)
	switch rt.Kind() {
	case reflect.Array:
		return d.arrayIter(rv)
	case reflect.Map:
		return d.mapIter(rt, rv)
	}
	return nil, nil
}

func (d *MapPart) Child(step Step) (Part, error) {
	i, err := d.Value()
	if err != nil {
		return nil, err
	}
	switch step.Type() {
	case TypeArray:
		ia, ok := i.([]interface{})
		if !ok {
			return nil, nil
		}
		key, _ := step.Int()
		if key < 0 || key >= len(ia) {
			return nil, nil
		}
		return NewMapPart(ia[key]), nil
	case TypeInt:
		im, ok := i.(map[int]interface{})
		if !ok {
			return nil, nil
		}
		i, ok := step.Int()
		if ok == false {
			return nil, nil
		}
		return NewMapPart(im[i]), nil
	case TypeString:
		sm, ok := i.(map[string]interface{})
		if !ok {
			return nil, nil
		}
		return NewMapPart(sm[step.String()]), nil
	}
	im, ok := i.(map[interface{}]interface{})
	if !ok {
		return nil, nil
	}
	return NewMapPart(im[step.Interface()]), nil
}

type mapElement struct {
	Step Step
	Part Part
}

type mapElements []mapElement

func (e *mapElements) Next() (*PartIter, error) {
	if len(*e) == 0 {
		return nil, nil
	}
	nme := mapElements((*e)[1:])
	return &PartIter{
		Step: (*e)[0].Step,
		Part: (*e)[0].Part,
		Next: nme.Next,
	}, nil
}
