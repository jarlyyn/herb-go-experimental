package assembler

import "reflect"

type TypeChecker struct {
	Type      interface{}
	CheckType func(a *Assembler) (bool, error)
}

func getReflectType(v interface{}) reflect.Type {
	return reflect.TypeOf(reflect.Indirect(reflect.ValueOf(v)))
}

var TypeCheckerString = &TypeChecker{
	Type: TypeString,
	CheckType: func(a *Assembler) (bool, error) {
		v, err := a.Value()
		if err != nil {
			return false, err
		}
		return getReflectType(v).Kind() == reflect.String, nil
	},
}

var TypeCheckerBool = &TypeChecker{
	Type: TypeBool,
	CheckType: func(a *Assembler) (bool, error) {
		v, err := a.Value()
		if err != nil {
			return false, err
		}
		return getReflectType(v).Kind() == reflect.Bool, nil
	},
}

var TypeCheckerInt = &TypeChecker{
	Type: TypeInt,
	CheckType: func(a *Assembler) (bool, error) {
		v, err := a.Value()
		if err != nil {
			return false, err
		}
		return getReflectType(v).Kind() == reflect.Int, nil
	},
}
var TypeCheckerUiit = &TypeChecker{
	Type: TypeUint,
	CheckType: func(a *Assembler) (bool, error) {
		v, err := a.Value()
		if err != nil {
			return false, err
		}
		return getReflectType(v).Kind() == reflect.Uint, nil
	},
}
var TypeCheckerInt64 = &TypeChecker{
	Type: TypeInt64,
	CheckType: func(a *Assembler) (bool, error) {
		v, err := a.Value()
		if err != nil {
			return false, err
		}
		return getReflectType(v).Kind() == reflect.Int64, nil
	},
}
var TypeCheckerUint64 = &TypeChecker{
	Type: TypeUint64,
	CheckType: func(a *Assembler) (bool, error) {
		v, err := a.Value()
		if err != nil {
			return false, err
		}
		return getReflectType(v).Kind() == reflect.Uint64, nil
	},
}
