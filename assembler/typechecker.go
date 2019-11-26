package assembler

import "reflect"

type TypeChecker struct {
	Type      Type
	CheckType func(a *Assembler, rt reflect.Type) (bool, error)
}

type TypeCheckers []*TypeChecker

func (c *TypeCheckers) CheckType(a *Assembler, rt reflect.Type) (Type, error) {
	for _, v := range *c {
		ok, err := v.CheckType(a, rt)
		if err != nil {
			return TypeUnkonwn, err
		}
		if ok {
			return v.Type, nil
		}
	}
	return TypeUnkonwn, nil
}
func (c *TypeCheckers) Append(checkers ...*TypeChecker) {
	*c = append(*c, checkers...)
}

func (c *TypeCheckers) Insert(checkers ...*TypeChecker) {
	*c = TypeCheckers(append(checkers, *c...))
}
func getReflectType(v interface{}) reflect.Type {
	return reflect.TypeOf(reflect.Indirect(reflect.ValueOf(v)))
}

var TypeCheckerString = &TypeChecker{
	Type: TypeString,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.String, nil
	},
}

var TypeCheckerBool = &TypeChecker{
	Type: TypeBool,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Bool, nil
	},
}

var TypeCheckerInt = &TypeChecker{
	Type: TypeInt,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Int, nil
	},
}

var TypeCheckerUint = &TypeChecker{
	Type: TypeUint,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Uint, nil
	},
}
var TypeCheckerInt64 = &TypeChecker{
	Type: TypeInt64,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Int64, nil
	},
}
var TypeCheckerUint64 = &TypeChecker{
	Type: TypeUint64,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Uint64, nil
	},
}
var TypeCheckerFloat32 = &TypeChecker{
	Type: TypeFloat32,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Float32, nil
	},
}
var TypeCheckerFloat64 = &TypeChecker{
	Type: TypeFloat64,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Float64, nil
	},
}
var TypeCheckerStringKeyMap = &TypeChecker{
	Type: TypeMap,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Map && rt.Key().Kind() == reflect.String, nil
	},
}

var TypeCheckerSlice = &TypeChecker{
	Type: TypeSlice,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Slice, nil
	},
}

var TypeCheckerStruct = &TypeChecker{
	Type: TypeStruct,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Struct, nil
	},
}

var TypeCheckerEmptyInterface = &TypeChecker{
	Type: TypeEmptyInterface,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Interface && rt.NumMethod() == 0, nil
	},
}

var TypeCheckerPtr = &TypeChecker{
	Type: TypePtr,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Ptr, nil
	},
}

func SetCommonTypeCheckers(c *TypeCheckers) {
	c.Append(
		TypeCheckerBool,
		TypeCheckerString,
		TypeCheckerInt,
		TypeCheckerUint,
		TypeCheckerInt64,
		TypeCheckerUint64,
		TypeCheckerFloat32,
		TypeCheckerFloat64,
		TypeCheckerStringKeyMap,
		TypeCheckerSlice,
		TypeCheckerStruct,
		TypeCheckerEmptyInterface,
		TypeCheckerPtr,
	)
}
