package unmarshaler

import "testing"

type testSubStruct struct {
	SubStructValue string
}

var testData = map[string]interface{}{
	"Value":              "value",
	"caseinsensitive":    "ci",
	"namedcasesensitive": "namedcs",
	"unexportedfield":    "unexportedfield",

	"FieldInt":    int(1),
	"FieldIntPtr": int(2),

	"FieldInt64":    int64(64),
	"FieldInt64Ptr": int64(65),

	"FieldFloat32":    float32(32.0),
	"FieldFloat32Ptr": float32(33.0),

	"FieldFloat64":    float32(64.0),
	"FieldFloat64Ptr": float32(65.0),

	"FieldBool":    true,
	"FieldBoolPtr": true,

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

	"FieldMap":   map[int]string{1: "1"},
	"FieldSlice": []string{"slice1"},

	"LazyLoadFunc":        "loaderfunc",
	"LazyLoader":          "loader",
	"AnonymousValue":      "AnonymousValueStr",
	"NamedAnonymousValue": "NamedAnonymousValueStr",
	"ExistAnonymous": map[string]interface{}{
		"ExistAnonymousValue": "ExistAnonymousValueStr",
	},
	"ciexistanonymous": map[string]interface{}{
		"ciexistanonymousValue": "CIExistAnonymousValueStr",
	},
	"SubStruct": testSubStruct{
		SubStructValue: "SubStructValueStr",
	},
}

type Anonymous struct {
	AnonymousValue string
}

type NamedAnonymous struct {
	NamedAnonymousValue string
}

type ExistAnonymous struct {
	ExistAnonymousValue string
}

type CIExistAnonymous struct {
	CIExistAnonymousValue string
}
type testStruct struct {
	NamedValue      string `config:"Value"`
	CaseInsensitive string
	CaseSensitive   string `config:"Namedcasesensitive"`
	unexportedfield string

	FieldInt         int
	FieldEmptyInt    int
	FieldIntPtr      *int
	FieldEmptyIntPtr *int

	FieldInt64         int64
	FieldEmptyInt64    int64
	FieldInt64Ptr      *int64
	FieldEmptyInt64Ptr *int64

	FieldFloat32         int64
	FieldEmptyFloat32    int64
	FieldFloat32Ptr      *int64
	FieldEmptyFloat32Ptr *int64

	FieldFloat64         int64
	FieldEmptyFloat64    int64
	FieldFloat64Ptr      *int64
	FieldEmptyFloat64Ptr *int64

	FieldInt64ToInt     int     `config:"FieldInt64"`
	FieldInt64ToFloat32 float32 `config:"FieldInt64"`
	FieldInt64ToFloat64 float64 `config:"FieldInt64"`

	FieldFloat64ToInt int `config:"FieldFloat64"`
	FieldFloat32ToInt int `config:"FieldFloat32"`

	FieldBool         bool
	FieldEmptyBool    bool
	FieldBoolPtr      *bool
	FieldEmptyBoolPtr *bool

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

	FieldMap   interface{}
	FieldSMap  interface{} `config:"FieldStringMap"`
	FieldSlice interface{}

	Anonymous
	NamedAnonymous `config:"NamedAnonymousNotExist"`
	ExistAnonymous
	CIExistAnonymous
	SubStruct testSubStruct
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

	if v.FieldFloat32 != 32 {
		t.Fatal(v)
	}
	if v.FieldEmptyFloat32 != 0 {
		t.Fatal(v)
	}
	if v.FieldFloat32Ptr == nil || *v.FieldFloat32Ptr != 33 {
		t.Fatal(v)
	}
	if v.FieldEmptyFloat32Ptr != nil {
		t.Fatal(v)
	}

	if v.FieldFloat64 != 64 {
		t.Fatal(v)
	}
	if v.FieldEmptyFloat64 != 0 {
		t.Fatal(v)
	}
	if v.FieldFloat64Ptr == nil || *v.FieldFloat64Ptr != 65 {
		t.Fatal(v)
	}
	if v.FieldEmptyFloat64Ptr != nil {
		t.Fatal(v)
	}

	if v.FieldBool != true {
		t.Fatal(v)
	}
	if v.FieldEmptyBool != false {
		t.Fatal(v)
	}
	if v.FieldBoolPtr == nil || *v.FieldBoolPtr != true {
		t.Fatal(v)
	}
	if v.FieldEmptyBoolPtr != nil {
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
	if v.unexportedfield != "" {
		t.Fatal(v)
	}

	if v.Anonymous.AnonymousValue != "AnonymousValueStr" {
		t.Fatal(v)
	}
	if v.NamedAnonymous.NamedAnonymousValue != "" {
		t.Fatal(v)
	}
	if v.ExistAnonymous.ExistAnonymousValue != "ExistAnonymousValueStr" {
		t.Fatal(v)
	}
	if v.CIExistAnonymous.CIExistAnonymousValue != "CIExistAnonymousValueStr" {
		t.Fatal(v)
	}
	m, ok := v.FieldMap.(map[interface{}]interface{})
	if !ok || m[1] != "1" {
		t.Fatal(v)
	}

	ms, ok := v.FieldSMap.(map[string]interface{})
	if !ok || ms["key1"] != "elem1" {
		t.Fatal(v)
	}
	s, ok := v.FieldSlice.([]interface{})
	if !ok || s[0] != "slice1" {
		t.Fatal(v)
	}
	if v.FieldInt64ToFloat64 != 64.0 {
		t.Fatal(v)
	}

	if v.FieldFloat64ToInt != 64 {
		t.Fatal(v)
	}

	if v.FieldFloat32ToInt != 32 {
		t.Fatal(v)
	}

	if v.SubStruct.SubStructValue != "SubStructValueStr" {
		t.Fatal(v)
	}

}
