package guarder

import (
	"net/http"
	"testing"
)

func TestIDTokenHeaders(t *testing.T) {
	req, err := http.NewRequest("POST", "http://127.0,0,1", nil)
	if err != nil {
		t.Fatal(err)
	}
	c := &ConfigMap{}
	c.Set("IDHeader", "id")
	d, err := NewMapperDriver("headers", c, "")
	if err != nil {
		t.Fatal(err)
	}
	p, err := d.ReadParamsFromRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if p.ID() != "" || p.Credential() != "" {
		t.Fatal(*p)
	}
	p = NewParams()
	p.SetID("testid")
	p.SetCredential("teestoken")
	d.WriteParamsToRequest(req, p)
	p, err = d.ReadParamsFromRequest(req)
	if err != nil {
		t.Fatal(err)
	}
	if p.ID() != "testid" || p.Credential() != "teestoken" {
		t.Fatal(*p)
	}
}
