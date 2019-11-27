package unmarshaler

type Type string

var TypeUnkonwn = Type("")
var TypeBool = Type("unmarshaler.bool")
var TypeString = Type("unmarshaler.string")
var TypeInt8 = Type("unmarshaler.int8")
var TypeUint8 = Type("unmarshaler.uint8")
var TypeInt16 = Type("unmarshaler.int16")
var TypeUint16 = Type("unmarshaler.uint16")
var TypeInt = Type("unmarshaler.int")
var TypeUint = Type("unmarshaler.uint")
var TypeInt64 = Type("unmarshaler.int64")
var TypeUint64 = Type("unmarshaler.uint64")
var TypeFloat32 = Type("unmarshaler.float32")
var TypeFloat64 = Type("unmarshaler.float64")

var TypeMap = Type("unmarshaler.map")
var TypeArray = Type("unmarshaler.array")
var TypeSlice = Type("unmarshaler.slice")
var TypeStruct = Type("unmarshaler.struct")
var TypeStructField = Type("unmarshaler.structFild")
var TypeEmptyInterface = Type("unmarshaler.interface{}")
var TypeLazyLoadFunc = Type("unmarshaler.lazyloadfunc")
var TypeLazyLoader = Type("unmarshaler.lazyloader")
var TypePtr = Type("unmarshaler.*")
