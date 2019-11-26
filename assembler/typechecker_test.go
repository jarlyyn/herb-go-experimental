package assembler

import "testing"

var testData = map[string]interface{}{
	"FieldInt": int(1),
}

type testStruct struct {
	FieldInt int
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
}
