package assembler

// type Type interface{}

// var TypeBool = Type(reflect.TypeOf(false))
// var ReflectTypeString = reflect.TypeOf("")
// var TypeString = Type(ReflectTypeString)
// var TypeInt8 = Type(reflect.TypeOf(int8(0)))
// var TypeUint8 = Type(reflect.TypeOf(uint8(0)))
// var TypeInt16 = Type(reflect.TypeOf(int16(0)))
// var TypeUint16 = Type(reflect.TypeOf(uint16(0)))
// var TypeInt = Type(reflect.TypeOf(int(0)))
// var TypeUint = Type(reflect.TypeOf(uint(0)))
// var TypeInt64 = Type(reflect.TypeOf(int64(0)))
// var TypeUint64 = Type(reflect.TypeOf(int64(0)))
// var TypeFloat32 = Type(reflect.TypeOf(float32(0)))
// var TypeFloat64 = Type(reflect.TypeOf(float64(0)))

// var TypeMap = Type(reflect.TypeOf(map[interface{}]interface{}{}))
// var TypeArray = Type(reflect.TypeOf([0]interface{}{}))
// var TypeSlice = Type(reflect.TypeOf([]interface{}{}))
// var TypeStruct = Type(reflect.TypeOf(&struct{}{}))
// var TypeStructField = Type(reflect.TypeOf(&reflect.StructField{}))
// var ReflectTypeEmptyInterface = reflect.TypeOf((*interface{})(nil)).Elem()
// var TypeEmptyInterface = Type(ReflectTypeEmptyInterface)

type Type string

var TypeUnkonwn = Type("")
var TypeBool = Type("bool")
var TypeString = Type("string")
var TypeInt8 = Type("int8")
var TypeUint8 = Type("uint8")
var TypeInt16 = Type("int16")
var TypeUint16 = Type("uint16")
var TypeInt = Type("int")
var TypeUint = Type("uint")
var TypeInt64 = Type("int64")
var TypeUint64 = Type("uint64")
var TypeFloat32 = Type("float32")
var TypeFloat64 = Type("float64")

var TypeMap = Type("map")
var TypeArray = Type("array")
var TypeSlice = Type("slice")
var TypeStruct = Type("struct")
var TypeStructField = Type("structFild")
var TypeEmptyInterface = Type("interface{}")
