package simpleapi

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/herb-go/util/httpserver"
)

func TestRouter(t *testing.T) {
	var err error
	Reset()
	CleanConfig()
	defer func() {
		Reset()
		CleanConfig()
	}()
	config := &httpserver.Config{
		Net:  "tcp",
		Addr: ":6789",
	}
	err = SetConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	Start("/test", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/test"))
	})
	defer Stop("/test")
	Start("/test2", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("/test2"))
	})
	defer Stop("/test2")
	resp, err := http.Get("http://" + config.Addr + "/test")
	if err != nil {
		t.Fatal(err)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(bs) != "/test" {
		t.Fatal(string(bs))
	}
	resp, err = http.Get("http://" + config.Addr + "/test2")
	if err != nil {
		t.Fatal(err)
	}
	bs, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if string(bs) != "/test2" {
		t.Fatal(string(bs))
	}
	resp, err = http.Get("http://" + config.Addr + "/notexist")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 404 {
		t.Fatal(resp)
	}
}
func TestDaemon(t *testing.T) {
	var err error
	defer func() {
		CleanConfig()
	}()
	if runningCount != 0 {
		t.Fatal(runningCount)
	}
	err = Stop("test")
	if err == nil {
		t.Fatal(err)
	}
	err = Start("test", func(w http.ResponseWriter, r *http.Request) {})
	if err != nil {
		t.Fatal(err)
	}
	if runningCount != 1 {
		t.Fatal(runningCount)
	}
	err = Start("test", func(w http.ResponseWriter, r *http.Request) {})
	if err == nil {
		t.Fatal(err)
	}
	if runningCount != 1 {
		t.Fatal(runningCount)
	}
	err = Stop("test")
	if err != nil {
		t.Fatal(err)
	}
	if runningCount != 0 {
		t.Fatal(runningCount)
	}
	err = Start("test", func(w http.ResponseWriter, r *http.Request) {})
	if err != nil {
		t.Fatal(err)
	}
	err = Start("test2", func(w http.ResponseWriter, r *http.Request) {})
	if err != nil {
		t.Fatal(err)
	}
	if runningCount != 2 {
		t.Fatal(runningCount)
	}
	err = Stop("test")
	if err != nil {
		t.Fatal(err)
	}
	if runningCount != 1 {
		t.Fatal(runningCount)
	}
	err = Stop("test2")
	if err != nil {
		t.Fatal(err)
	}
	if runningCount != 0 {
		t.Fatal(runningCount)
	}
	err = Stop("test")
	if err == nil {
		t.Fatal(err)
	}
	if runningCount != 0 {
		t.Fatal(runningCount)
	}
}
