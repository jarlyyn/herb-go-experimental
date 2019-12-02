package unmarshaler

import (
	"reflect"
)

//MapPart assembler interface map part struct
type MapPart struct {
	value interface{}
}

//NewMapPart create new map part
func NewMapPart(v interface{}) *MapPart {
	return &MapPart{
		value: v,
	}
}

//Value return part value as empty interface.
//Shold only used when part is not iterable
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

//Iter return part iter.
//Nil should be returned If part is not iterable.
func (d *MapPart) Iter() (*PartIter, error) {
	v, err := d.Value()
	if err != nil {
		return nil, err
	}
	rt := reflect.TypeOf(v)
	rv := reflect.ValueOf(v)
	switch rt.Kind() {
	case reflect.Array, reflect.Slice:
		return d.arrayIter(rv)
	case reflect.Map:
		return d.mapIter(rt, rv)
	}
	return nil, nil
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
