package query

import (
	"strings"
)

func NewValueList(data ...interface{}) *PlainQuery {
	if len(data) == 0 {
		return New("")
	}
	var command = strings.Repeat("? , ", len(data))
	return New(command[:len(command)-3], data...)
}

func In(data ...interface{}) *PlainQuery {
	var query = NewValueList(data...)
	query.Command = "IN ( " + query.Command + " )"
	return query
}
