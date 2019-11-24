package assembler

import (
	"reflect"
)

type CommonType interface{}

var TypeBool = CommonType(reflect.TypeOf(false))
var TypeString = CommonType(reflect.TypeOf(""))
var TypeInt = CommonType(reflect.TypeOf(int(0)))
var TypeUint = CommonType(reflect.TypeOf(uint(0)))
var TypeInt64 = CommonType(reflect.TypeOf(int64(0)))
var TypeUint64 = CommonType(reflect.TypeOf(int64(0)))

var TypeInterface = CommonType(reflect.TypeOf((*interface{})(nil)).Elem())
var TypeMap = CommonType(reflect.TypeOf(map[interface{}]interface{}{}))
var TypeArray = CommonType(reflect.TypeOf([]interface{}{}))
var TypeStruct = CommonType(reflect.TypeOf(struct{}{}))
