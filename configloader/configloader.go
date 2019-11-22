package configloader

import (
	"reflect"
)

type ConfigDecoder interface {
	DecodeConfig(interface{}) error
}

var rtConfigDecoder = reflect.TypeOf((*ConfigDecoder)(nil)).Elem()

func IsConfigDecoder(rt reflect.Type) bool {
	return rt.Implements(rtConfigDecoder)
}
