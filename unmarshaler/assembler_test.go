package unmarshaler

import "testing"

var testData = map[string]interface{}{
	"Value":              "value",
	"caseinsensitive":    "ci",
	"namedcasesensitive": "namedcs",
	"FieldInt":           int(1),
	"FieldIntPtr":        int(2),

	"FieldInt64":    int64(64),
	"FieldInt64Ptr": int64(65),

	"FieldString":    "str",
	"FieldStringPtr": "str2",

	"FieldStringSlice":    []string{"elem1", "elem2"},
	"FieldStringSlicePtr": []string{"elem3", "elem4"},

	"FieldInterfaceSlice":    []interface{}{"elemi1", "elemi2"},
	"FieldInterfaceSlicePtr": []interface{}{"elemi3", "elemi4"},

	"FieldStringMap":    map[string]interface{}{"key1": "elem1", "key2": "elem2"},
	"FieldStringMapPtr": map[string]interface{}{"key3": "elem3", "key4": "elem4"},

	"FieldInterfaceMap":    map[string]interface{}{"key1": "elem1", "key2": "elem2"},
	"FieldInterfaceMapPtr": map[string]interface{}{"key3": "elem3", "key4": "elem4"},

	"LazyLoadFunc": "loaderfunc",
	"LazyLoader":   "loader",
}

type testStruct struct {
	NamedValue       string `config:"Value"`
	CaseInsensitive  string
	CaseSensitive    string `config:"Namedcasesensitive"`
	FieldInt         int
	FieldEmptyInt    int
	FieldIntPtr      *int
	FieldEmptyIntPtr *int

	FieldInt64         int64
	FieldEmptyInt64    int64
	FieldInt64Ptr      *int64
	FieldEmptyInt64Ptr *int64

	FieldString         string
	FieldEmptyString    string
	FieldStringPtr      *string
	FieldEmptyStringPtr *string

	FieldStringSlice         []string
	FieldEmptyStringSlice    []string
	FieldStringSlicePtr      *[]string
	FieldEmptyStringSlicePtr *[]string

	FieldInterfaceSlice         []interface{}
	FieldEmptyInterfaceSlice    []interface{}
	FieldInterfaceSlicePtr      *[]interface{}
	FieldEmptyInterfaceSlicePtr *[]interface{}

	FieldInterfaceMap         map[string]interface{}
	FieldEmptyInterfaceMap    map[string]interface{}
	FieldInterfaceMapPtr      *map[string]interface{}
	FieldEmptyInterfaceMapPtr *map[string]interface{}

	FieldStringMap         map[string]interface{}
	FieldEmptyStringMap    map[string]interface{}
	FieldStringMapPtr      *map[string]interface{}
	FieldEmptyStringMapPtr *map[string]interface{}
	LazyLoadFunc           func(v interface{}) error `config:", lazyload"`
	NamedLazyLoader        LazyLoader                `config:"LazyLoader,lazyload"`
}

func TestAssembler(t *testing.T) {
	c := NewCommonConfig()
	a := EmptyAssembler.WithConfig(c).WithPart(NewMapPart(testData))
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

	if v.FieldInt64 != 64 {
		t.Fatal(v)
	}
	if v.FieldEmptyInt64 != 0 {
		t.Fatal(v)
	}
	if v.FieldInt64Ptr == nil || *v.FieldInt64Ptr != 65 {
		t.Fatal(v)
	}
	if v.FieldEmptyInt64Ptr != nil {
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
	if len(v.FieldInterfaceSlice) != 2 || v.FieldInterfaceSlice[0] != "elemi1" || v.FieldInterfaceSlice[1] != "elemi2" {
		t.Fatal(v)
	}

	if len(v.FieldEmptyInterfaceSlice) != 0 {
		t.Fatal(v)
	}
	if v.FieldInterfaceSlicePtr == nil || len(*v.FieldInterfaceSlicePtr) != 2 || (*v.FieldInterfaceSlicePtr)[0] != "elemi3" || (*v.FieldInterfaceSlicePtr)[1] != "elemi4" {
		t.Fatal(v)
	}

	if v.FieldEmptyInterfaceSlicePtr != nil {
		t.Fatal(v)
	}

	if len(v.FieldInterfaceMap) != 2 || v.FieldInterfaceMap["key1"] != "elem1" || v.FieldInterfaceMap["key2"] != "elem2" {
		t.Fatal(v)
	}

	if len(v.FieldEmptyInterfaceMap) != 0 {
		t.Fatal(v)
	}
	if v.FieldInterfaceMapPtr == nil || len(*v.FieldInterfaceMapPtr) != 2 || (*v.FieldInterfaceMapPtr)["key3"] != "elem3" || (*v.FieldInterfaceMapPtr)["key4"] != "elem4" {
		t.Fatal(v)
	}

	if v.FieldEmptyInterfaceMapPtr != nil {
		t.Fatal(v)
	}
	lazyloaderesult := ""
	if v.LazyLoadFunc == nil {
		t.Fatal(v)
	}
	if err := v.LazyLoadFunc(&lazyloaderesult); err != nil {
		t.Fatal(err)
	}
	if lazyloaderesult != "loaderfunc" {
		t.Fatal(lazyloaderesult)
	}
	lazyloaderesult = ""
	if v.NamedLazyLoader == nil {
		t.Fatal(v)
	}
	if err := v.NamedLazyLoader.LazyLoad(&lazyloaderesult); err != nil {
		t.Fatal(err)
	}
	if lazyloaderesult != "loader" {
		t.Fatal(lazyloaderesult)
	}
	if v.NamedValue != "value" {
		t.Fatal(v)
	}
	if v.CaseInsensitive != "ci" {
		t.Fatal(v)
	}
	if v.CaseSensitive != "" {
		t.Fatal(v)
	}
}
