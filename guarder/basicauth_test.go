package guarder

import (
	"net/http"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	req, err := http.NewRequest("POST", "http://127.0,0,1", nil)
	if err != nil {
		t.Fatal(err)
	}
	d, err := NewMapperDriver("basicauth", &ConfigMap{}, "")
	if err != nil {
		t.Fatal(err)
	}
	p, err := d.ReadParamsFromRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if p.ID() != "" || p.Token() != "" {
		t.Fatal(*p)
	}
	p = NewParams()
	p.SetID("testid")
	p.SetToken("teestoken")
	d.WriteParamsToRequest(req, p)
	p, err = d.ReadParamsFromRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if p.ID() != "testid" || p.Token() != "teestoken" {
		t.Fatal(*p)
	}
}
