package herbunit

import "reflect"

type Unit struct {
	Path      string
	Name      string
	Component interface{}
}

func New() *Unit {
	return &Unit{}
}

type UnitType interface {
	InitUnitType() error
	Summary() (interface{}, error)
	PlainSummary() (string, error)
	Keyword() string
}

var units = map[string][]*Unit{}
var types = map[string]UnitType{}

func Register(path string, name string, component interface{}) *Unit {
	ct := reflect.ValueOf(component).Elem().Type().String()
	if units[ct] == nil {
		units[ct] = []*Unit{}
	}
	c := New()
	c.Path = path
	c.Name = name
	c.Component = component
	units[ct] = append(units[ct], c)
	return c
}
func RegisterType(t UnitType, component interface{}) {
	types[reflect.ValueOf(component).Elem().Type().String()] = t
}

func Init() error {
	for k := range types {
		err := types[k].InitUnitType()
		if err != nil {
			return err
		}
	}
	return nil
}

func Reset() {
	units = map[string][]*Unit{}
	types = map[string]UnitType{}
}
