package configloader

import (
	"reflect"
	"testing"
)

type testInterface interface {
	DecodeConfig(interface{}) error
}

type testStruct struct {
	I testInterface
}

type testDecode struct {
}

func (d *testDecode) DecodeConfig(interface{}) error {
	return nil
}

type testStruct2 struct {
	I testDecode
}

func TestConfigDecoder(t *testing.T) {
	ts := testStruct{}
	tif, _ := reflect.TypeOf(ts).FieldByName("I")
	if !IsConfigDecoder(tif.Type) {
		t.Fatal(tif)
	}
	ts2 := testStruct{}
	tif2, _ := reflect.TypeOf(ts2).FieldByName("I")
	if IsConfigDecoder(tif2.Type) {
		t.Fatal(tif2)
	}

}
