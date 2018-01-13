package query

import (
	"reflect"
	"strings"
)

func NewValueList(data ...interface{}) *PlainQuery {
	if len(data) == 0 {
		return New("")
	}
	var command = strings.Repeat("? , ", len(data))
	return New(command[:len(command)-3], data...)
}

func In(field string, args interface{}) *PlainQuery {
	var argsvalue = reflect.ValueOf(args)
	var data = make([]interface{}, argsvalue.Len())
	for k := range data {
		data[k] = argsvalue.Index(k).Interface()
	}
	var query = NewValueList(data...)
	query.Command = field + " IN ( " + query.Command + " )"
	return query
}
