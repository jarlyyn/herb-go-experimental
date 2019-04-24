package routeridentifier

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type stuctTestIdentifier struct {
	TagFields map[string]bool
}

func (i *stuctTestIdentifier) MustIdentifyRouter(prefix string, r *http.Request) {
	id := NewIdentification()
	id.ID = prefix + r.URL.Path
	for k := range i.TagFields {
		if i.TagFields[k] {
			if r.Header.Get(k) != "" {
				id.AddTag(r.Header.Get(k))
			}
		}
	}
	SetIdentificationToRequest(r, id)
}

var testIdentifier = &stuctTestIdentifier{
	TagFields: map[string]bool{
		"testheader1": true,
		"testheader2": true,
	},
}

func TestIdentification(t *testing.T) {
	id := NewIdentification()
	id.AddTag("")
	id.AddTag("test")
	id.AddTag("test2")
	id.AddTag("test2")
	if !id.HasTag("test") {
		t.Fatal(id)
	}
	if !id.HasTag("test2") {
		t.Fatal(id)
	}
	if id.HasTag("test3") {
		t.Fatal(id)
	}
	id = NewIdentification()
	id.ID = "test"
	if id.String() != "test" {
		t.Fatal(id)
	}
}
func TestIdentifier(t *testing.T) {
	Debug = false

	var testAction = func(w http.ResponseWriter, r *http.Request) {
		id := GetIdentificationFromRequest(r)
		if id == nil {
			w.Write([]byte(""))
			return
		}
		bs, err := json.Marshal(id)
		if err != nil {
			t.Fatal(err)
		}
		w.Write(bs)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		Middleware(testIdentifier, "prefix")(w, r, testAction)
	})
	mux.HandleFunc("/nil", testAction)
	server := httptest.NewServer(mux)
	defer server.Close()
	testreq, err := http.NewRequest("GET", server.URL+"/test", nil)
	testreq.Header.Add("testheader1", "testvalue1")
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(testreq)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	result := NewIdentification()
	err = json.Unmarshal(content, result)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "prefix/test" {
		t.Fatal(result)
	}
	if !result.HasTag("testvalue1") {
		t.Fatal(result)
	}
	if resp.Header.Get(DebugHeader) != "" {
		t.Fatal(resp)
	}
	Debug = true
	defer func() {
		Debug = false
	}()
	testreq, err = http.NewRequest("GET", server.URL+"/test", nil)
	testreq.Header.Add("testheader1", "testvalue1")
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(testreq)
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Header.Get(DebugHeader) == "" {
		t.Fatal(resp)
	}
	resp.Body.Close()
	testreq, err = http.NewRequest("GET", server.URL+"/nil", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(testreq)
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if len(content) != 0 {
		t.Fatal(string(content))
	}
}
