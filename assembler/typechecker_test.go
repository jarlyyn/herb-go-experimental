package assembler

import "testing"

var testData = map[string]interface{}{
	"FieldInt":               int(1),
	"FieldIntPtr":            int(2),
	"FieldString":            "str",
	"FieldStringPtr":         "str2",
	"FieldStringSlice":       []string{"elem1", "elem2"},
	"FieldStringSlicePtr":    []string{"elem3", "elem4"},
	"FieldInterfaceSlice":    []interface{}{"elem1", "elem2"},
	"FieldInterfaceSlicePtr": []interface{}{"elem3", "elem4"},
}

type testStruct struct {
	FieldInt                    int
	FieldEmptyInt               int
	FieldIntPtr                 *int
	FieldEmptyIntPtr            *int
	FieldString                 string
	FieldEmptyString            string
	FieldStringPtr              *string
	FieldEmptyStringPtr         *string
	FieldStringSlice            []string
	FieldEmptyStringSlice       []string
	FieldStringSlicePtr         *[]string
	FieldEmptyStringSlicePtr    *[]string
	FieldInterfaceSlice         []interface{}
	FieldEmptyInterfaceSlice    []interface{}
	FieldInterfaceSlicePtr      *[]interface{}
	FieldEmptyInterfaceSlicePtr *[]interface{}
}

func TestTypeCheck(t *testing.T) {
	c := NewCommonConfig()
	a := RootAssembler.WithConfig(c).WithPart(NewMapPart(testData))
	v := &testStruct{}
	ok, err := a.Assemble(v)
	if ok == false {
		t.Fatal(ok)
	}
	if err != nil {
		t.Fatal(err)
	}
	if v.FieldInt != 1 {
		t.Fatal(v)
	}
	if v.FieldEmptyInt != 0 {
		t.Fatal(v)
	}
	if v.FieldIntPtr == nil || *v.FieldIntPtr != 2 {
		t.Fatal(v)
	}
	if v.FieldEmptyIntPtr != nil {
		t.Fatal(v)
	}
	if v.FieldString != "str" {
		t.Fatal(v)
	}
	if v.FieldEmptyString != "" {
		t.Fatal(v)
	}
	if v.FieldStringPtr == nil || *v.FieldStringPtr != "str2" {
		t.Fatal(v)
	}
	if v.FieldEmptyStringPtr != nil {
		t.Fatal(v)
	}
	if len(v.FieldStringSlice) != 2 || v.FieldStringSlice[0] != "elem1" || v.FieldStringSlice[1] != "elem2" {
		t.Fatal(v)
	}

	if len(v.FieldEmptyStringSlice) != 0 {
		t.Fatal(v)
	}
	if v.FieldStringSlicePtr == nil || len(*v.FieldStringSlicePtr) != 2 || (*v.FieldStringSlicePtr)[0] != "elem3" || (*v.FieldStringSlicePtr)[1] != "elem4" {
		t.Fatal(v)
	}

	if v.FieldEmptyStringSlicePtr != nil {
		t.Fatal(v)
	}
	if len(v.FieldInterfaceSlice) != 2 || v.FieldInterfaceSlice[0] != "elem1" || v.FieldInterfaceSlice[1] != "elem2" {
		t.Fatal(v)
	}

	if len(v.FieldEmptyInterfaceSlice) != 0 {
		t.Fatal(v)
	}
	if v.FieldInterfaceSlicePtr == nil || len(*v.FieldInterfaceSlicePtr) != 2 || (*v.FieldInterfaceSlicePtr)[0] != "elem3" || (*v.FieldStringSlicePtr)[1] != "elem4" {
		t.Fatal(v)
	}

	if v.FieldEmptyInterfaceSlicePtr != nil {
		t.Fatal(v)
	}
}
