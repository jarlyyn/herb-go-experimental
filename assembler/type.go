package assembler

import (
	"reflect"
)

type CommonType interface{}

var TypeBool = CommonType(reflect.TypeOf(false))
var TypeString = CommonType(reflect.TypeOf(""))
var TypeInt8 = CommonType(reflect.TypeOf(int8(0)))
var TypeUint8 = CommonType(reflect.TypeOf(uint8(0)))
var TypeInt16 = CommonType(reflect.TypeOf(int16(0)))
var TypeUint16 = CommonType(reflect.TypeOf(uint16(0)))
var TypeInt = CommonType(reflect.TypeOf(int(0)))
var TypeUint = CommonType(reflect.TypeOf(uint(0)))
var TypeInt64 = CommonType(reflect.TypeOf(int64(0)))
var TypeUint64 = CommonType(reflect.TypeOf(int64(0)))
var TypeFloat32 = CommonType(reflect.TypeOf(float32(0)))
var TypeFloat64 = CommonType(reflect.TypeOf(float64(0)))

var TypeMap = CommonType(reflect.TypeOf(map[interface{}]interface{}{}))
var TypeArray = CommonType(reflect.TypeOf([0]interface{}{}))
var TypeStruct = CommonType(reflect.TypeOf(&struct{}{}))
var TypeStructField = CommonType(reflect.TypeOf(&reflect.StructField{}))
var TypeInterface = CommonType(reflect.TypeOf((*interface{})(nil)).Elem())
