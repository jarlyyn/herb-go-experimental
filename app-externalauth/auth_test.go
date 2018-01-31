package auth

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/herb-go/herb/cache-session"
)

func TestServer(t *testing.T) {
	Auth := New()

	var successAction = func(w http.ResponseWriter, r *http.Request) {
		result := Auth.MustGetResult(r)
		data, err := json.Marshal(result)
		if err != nil {
			panic(err)
		}
		w.Write(data)
	}
	Session := session.NewClientStore([]byte("test"), 1*time.Hour)

	mux := http.StripPrefix("/auth", http.HandlerFunc(Auth.Serve(successAction)))
	server := httptest.NewServer(mux)
	Auth.MustInit(server.URL+"/auth", Session)
	data := Profile{
		ProfileIndexID: []string{"test"},
	}
	Auth.MustRegisterProvider("test", newTestDriver(data))
	provider := Auth.MustGetProvider("test")
	resp, err := http.Get(provider.LoginUrl())
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	result := &Result{}
	err = json.Unmarshal(content, result)
	if err != nil {
		t.Fatal(err)
	}
	if result.Account != data[ProfileIndexID][0] {
		t.Error(result.Account)
	}
	if result.Data[ProfileIndexID][0] != data[ProfileIndexID][0] {
		t.Error(result.Data[ProfileIndexID][0])
	}
}
