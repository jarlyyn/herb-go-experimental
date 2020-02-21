package responsecache

import (
	"net/http"
	"testing"

	"github.com/herb-go/herb/cache"
)

func TestParams(t *testing.T) {
	r, _ := http.NewRequest("GET", "http://127.0.0.1?query=test", nil)
	fp := FixedParam("fixed")
	qp := QueryParam("query")
	mp := MethodParam
	failp := ParamFunc(func(*http.Request) (string, bool) {
		return "fail", false
	})
	cb := NewParamsContextBuilder()
	fcb := cb.AppendParams(fp)
	if len(fcb.params) != 1 {
		t.Fatal(fcb)
	}
	if fcb.Identifier(r) != "fixed" {
		t.Fatal(fcb)
	}
	fqcb := fcb.AppendParams(qp)
	if len(fqcb.params) != 2 {
		t.Fatal(fqcb)
	}
	if fqcb.Identifier(r) != "fixed"+cache.KeyPrefix+"test" {
		t.Fatal(fqcb)
	}
	mfqcb := fqcb.AppendParams(mp)
	if len(mfqcb.params) != 3 {
		t.Fatal(mfqcb)
	}
	if mfqcb.Identifier(r) != "fixed"+cache.KeyPrefix+"test"+cache.KeyPrefix+"GET" {
		t.Fatal(mfqcb)
	}
	failcb := mfqcb.AppendParams(failp)
	if len(failcb.params) != 4 {
		t.Fatal(mfqcb)
	}
	if failcb.Identifier(r) != "" {
		t.Fatal(failcb)
	}

}
