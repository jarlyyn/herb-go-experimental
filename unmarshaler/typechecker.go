package unmarshaler

import "reflect"

//TypeChecker value type checker struct
type TypeChecker struct {
	Type      Type
	CheckType func(a *Assembler, rt reflect.Type) (bool, error)
}

//TypeCheckers type checkers list in order type
type TypeCheckers []*TypeChecker

//Append append checkers to last of given type checker.
func (c *TypeCheckers) Append(checkers ...*TypeChecker) *TypeCheckers {
	*c = append(*c, checkers...)
	return c
}

//AppendWith append with given TypeCheckers
func (c *TypeCheckers) AppendWith(checkers *TypeCheckers) *TypeCheckers {
	return c.Append(*checkers...)
}

//Insert insert checkers to first of given type checker.
func (c *TypeCheckers) Insert(checkers ...*TypeChecker) *TypeCheckers {
	*c = TypeCheckers(append(checkers, *c...))
	return c
}

//InsertWith insert with given TypeCheckers
func (c *TypeCheckers) InsertWith(checkers *TypeCheckers) *TypeCheckers {
	return c.Insert(*checkers...)
}

//NewTypeCheckers create new type checkers
func NewTypeCheckers() *TypeCheckers {
	return &TypeCheckers{}
}

//TypeCheckerString type checker for string.
var TypeCheckerString = &TypeChecker{
	Type: TypeString,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.String, nil
	},
}

//TypeCheckerBool type checker for bool.
var TypeCheckerBool = &TypeChecker{
	Type: TypeBool,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Bool, nil
	},
}

//TypeCheckerInt type checker for int.
var TypeCheckerInt = &TypeChecker{
	Type: TypeInt,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Int, nil
	},
}

//TypeCheckerUint type checker for uint.
var TypeCheckerUint = &TypeChecker{
	Type: TypeUint,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Uint, nil
	},
}

//TypeCheckerInt64 type checker for int64
var TypeCheckerInt64 = &TypeChecker{
	Type: TypeInt64,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Int64, nil
	},
}

//TypeCheckerUint64 type checker for uint64
var TypeCheckerUint64 = &TypeChecker{
	Type: TypeUint64,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Uint64, nil
	},
}

//TypeCheckerFloat32 type checker for float32
var TypeCheckerFloat32 = &TypeChecker{
	Type: TypeFloat32,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Float32, nil
	},
}

//TypeCheckerFloat64 type checker for float64
var TypeCheckerFloat64 = &TypeChecker{
	Type: TypeFloat64,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Float64, nil
	},
}

//TypeCheckerStringKeyMap type checker for string key map.
var TypeCheckerStringKeyMap = &TypeChecker{
	Type: TypeMap,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Map && rt.Key().Kind() == reflect.String, nil
	},
}

//TypeCheckerSlice type checker for slice
var TypeCheckerSlice = &TypeChecker{
	Type: TypeSlice,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Slice, nil
	},
}

//TypeCheckerStruct type checker for struct
var TypeCheckerStruct = &TypeChecker{
	Type: TypeStruct,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Struct, nil
	},
}

//TypeCheckerEmptyInterface type checker for empty interface.
var TypeCheckerEmptyInterface = &TypeChecker{
	Type: TypeEmptyInterface,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Interface && rt.NumMethod() == 0, nil
	},
}

//TypeCheckerLazyLoadFunc type checker for lazy load func.
var TypeCheckerLazyLoadFunc = &TypeChecker{
	Type: TypeLazyLoadFunc,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		lt := a.Config().TagLazyLoad
		if lt == "" {
			return false, nil
		}
		step := a.Step()
		if step == nil || step.Type() != TypeStructField {
			return false, nil
		}
		field := step.Interface().(reflect.StructField)
		tag, err := a.Config().GetTag(rt, field)
		if err != nil {
			return false, err
		}
		return rt.Kind() == reflect.Func && tag != nil && tag.Flags[lt] != "", nil
	},
}

//TypeCheckerLazyLoader type checker for lazy loader.
var TypeCheckerLazyLoader = &TypeChecker{
	Type: TypeLazyLoader,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		lt := a.Config().TagLazyLoad
		if lt == "" {
			return false, nil
		}
		step := a.Step()
		if step == nil || step.Type() != TypeStructField {
			return false, nil
		}
		field := step.Interface().(reflect.StructField)
		tag, err := a.Config().GetTag(rt, field)
		if err != nil {
			return false, err
		}
		return rt.Kind() == reflect.Interface && tag != nil && tag.Flags[lt] != "", nil
	},
}

//TypeCheckerPtr type checker for pointer
var TypeCheckerPtr = &TypeChecker{
	Type: TypePtr,
	CheckType: func(a *Assembler, rt reflect.Type) (bool, error) {
		return rt.Kind() == reflect.Ptr, nil
	},
}

//DefaultCommonTypeCheckers default common type checkers
func DefaultCommonTypeCheckers() *TypeCheckers {
	return NewTypeCheckers().Append(
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
		TypeCheckerLazyLoadFunc,
		TypeCheckerLazyLoader,
	)
}

//CommonTypeCheckers common type checkers used in NewCommonConfig
var CommonTypeCheckers = NewTypeCheckers().AppendWith(DefaultCommonTypeCheckers())
