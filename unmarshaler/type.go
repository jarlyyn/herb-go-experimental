package unmarshaler

//Type field type
type Type string

//TypeUnkonwn field type unkowwn
var TypeUnkonwn = Type("")

//TypeBool field type bool
var TypeBool = Type("unmarshaler.bool")

//TypeString field type string
var TypeString = Type("unmarshaler.string")

//TypeInt8 field type int8
var TypeInt8 = Type("unmarshaler.int8")

//TypeUint8 field type uint8
var TypeUint8 = Type("unmarshaler.uint8")

//TypeInt16 field type int16
var TypeInt16 = Type("unmarshaler.int16")

//TypeUint16 field type uint16
var TypeUint16 = Type("unmarshaler.uint16")

//TypeInt filed type int
var TypeInt = Type("unmarshaler.int")

//TypeUint field type uint
var TypeUint = Type("unmarshaler.uint")

//TypeInt64 field type int64
var TypeInt64 = Type("unmarshaler.int64")

//TypeUint64 field type uint64
var TypeUint64 = Type("unmarshaler.uint64")

//TypeFloat32 field type float32
var TypeFloat32 = Type("unmarshaler.float32")

//TypeFloat64 field type float64
var TypeFloat64 = Type("unmarshaler.float64")

//TypeMap field type map
var TypeMap = Type("unmarshaler.map")

//TypeArray field type array
var TypeArray = Type("unmarshaler.array")

//TypeSlice field type slice
var TypeSlice = Type("unmarshaler.slice")

//TypeStruct field type struct
var TypeStruct = Type("unmarshaler.struct")

//TypeStructField field type struct field
var TypeStructField = Type("unmarshaler.structFild")

//TypeEmptyInterface field type empty interface
var TypeEmptyInterface = Type("unmarshaler.interface{}")

//TypeLazyLoadFunc field type lazyload func
var TypeLazyLoadFunc = Type("unmarshaler.lazyloadfunc")

//TypeLazyLoader field type lazyloader
var TypeLazyLoader = Type("unmarshaler.lazyloader")

//TypePtr field type pointer
var TypePtr = Type("unmarshaler.*")
