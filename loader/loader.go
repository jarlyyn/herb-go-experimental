package main

import (
	"encoding"
	"fmt"
	"reflect"
)

type loader interface {
	UnmarshalBinary(data []byte) error
}

var l = (*loader)(nil)
var loadertype = reflect.TypeOf(l).Elem()

type teststruct struct {
	Data encoding.BinaryUnmarshaler
}

type u struct {
}

func (u *u) UnmarshalBinary(data []byte) error {
	return nil
}
func main() {
	var test = &teststruct{
		Data: &u{},
	}
	t := reflect.TypeOf(test).Elem()
	fmt.Println(t)
	df, _ := t.FieldByName("Data")
	v := reflect.ValueOf(test)
	tv := v.Type().Elem()
	dfv, _ := tv.FieldByName("Data")
	fmt.Println(df)
	fmt.Println(df.Type)
	fmt.Println(dfv)
	fmt.Println(dfv.Type)
	fmt.Println(loadertype)
	fmt.Println(df.Type.Implements(loadertype))
	fmt.Println(dfv.Type.Implements(loadertype))

	// fmt.Println(t.Implements(loadertype))
}
