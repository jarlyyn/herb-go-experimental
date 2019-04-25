package httprouteridentifier

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"

	"github.com/jarlyyn/herb-go-experimental/pathid"
)

func TestIdentifier(t *testing.T) {
	var err error
	var req *http.Request
	var resp *http.Response
	var content []byte
	var result *pathid.Identification
	idfer := NewIndentifier()
	config := NewConfig()
	config.Enabled = true
	var ActionNormal = NewAction()
	ActionNormal.Enabled = true
	ActionNormal.Path = "/normal"
	ActionNormal.ID = "normal"
	ActionNormal.Method = "GET"
	ActionNormal.Tags = []string{"normaltag"}
	config.Router = append(config.Router, ActionNormal)
	var ActionDisabled = NewAction()
	ActionDisabled.Enabled = false
	ActionDisabled.Path = "/disabled"
	ActionDisabled.ID = "disabled"
	ActionDisabled.Method = "ALL"
	ActionDisabled.Tags = []string{"disabled"}
	config.Router = append(config.Router, ActionDisabled)
	var ActionSubRouter = NewAction()
	ActionSubRouter.Enabled = true
	ActionSubRouter.Path = "/subrouter/:path"
	ActionSubRouter.ID = "subrouterid"
	ActionSubRouter.IsSubRouter = true
	ActionSubRouter.Method = ""
	ActionSubRouter.Tags = []string{"subroutertag"}
	config.Router = append(config.Router, ActionSubRouter)
	idfer.SubRouters["subrouterid"] = httprouter.New()
	config.SubRouters["subrouterid"] = NewRouter()
	var ActionSubRouterNormal = NewAction()
	ActionSubRouterNormal.Enabled = true
	ActionSubRouterNormal.Path = "/normal"
	ActionSubRouterNormal.ID = "subrouternormal"
	ActionSubRouterNormal.Method = "POST"
	ActionSubRouterNormal.Tags = []string{"subrouternormaltag"}
	config.SubRouters["subrouterid"] = append(config.SubRouters["subrouterid"], ActionSubRouterNormal)
	err = config.ApplyTo(idfer)
	if err != nil {
		t.Fatal(err)
	}
	var testAction = func(w http.ResponseWriter, r *http.Request) {
		id := pathid.GetIdentificationFromRequest(r)
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathid.Middleware(idfer, "testhost")(w, r, testAction)
	}))
	defer server.Close()
	req, err = http.NewRequest("GET", server.URL+"/normal", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	result = pathid.NewIdentification()
	err = json.Unmarshal(content, result)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "testhost/normal#GET" || !result.HasTag("normaltag") {
		t.Fatal(result)
	}
	req, err = http.NewRequest("POST", server.URL+"/subrouter"+"/normal", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	result = pathid.NewIdentification()
	err = json.Unmarshal(content, result)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "testhost/subrouternormal#POST" || !result.HasTag("subroutertag") || !result.HasTag("subrouternormaltag") || len(result.Parents) != 1 || result.Parents[0] != "subrouterid" {
		t.Fatal(result)
	}
	req, err = http.NewRequest("POST", server.URL+"/subrouter"+"/notexist", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	result = pathid.NewIdentification()
	err = json.Unmarshal(content, result)
	if err != nil {
		t.Fatal(err)
	}
	if result.ID != "" || !result.HasTag("subroutertag") {
		t.Fatal(result)
	}
	req, err = http.NewRequest("POST", server.URL+"/disabled", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	content, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
	if len(content) != 0 {
		t.Fatal(content)
	}
}

func TestActions(t *testing.T) {
	var err error
	idfer := NewIndentifier()
	var ActionWrongPath = NewAction()
	ActionWrongPath.Enabled = true
	ActionWrongPath.Path = "/wrongpath"
	ActionWrongPath.IsSubRouter = true
	ActionWrongPath.ID = "wrongpath"
	ActionWrongPath.Method = "GET"
	err = ActionWrongPath.ApplyTo(idfer, idfer.Router)
	if err == nil {
		t.Fatal(err)
	}

	var ActionWrongMethod = NewAction()
	ActionWrongMethod.Enabled = true
	ActionWrongMethod.Path = "/wrongmethod"
	ActionWrongMethod.ID = "wrongmethod"
	ActionWrongMethod.Method = "asd"
	err = ActionWrongMethod.ApplyTo(idfer, idfer.Router)
	if err == nil {
		t.Fatal(err)
	}
	var ActionMethodAutoUpper = NewAction()
	ActionMethodAutoUpper.Enabled = true
	ActionMethodAutoUpper.Path = "/autoapper"
	ActionMethodAutoUpper.ID = "autoapper"
	ActionMethodAutoUpper.Method = "Get"
	err = ActionMethodAutoUpper.ApplyTo(idfer, idfer.Router)
	if err != nil {
		t.Fatal(err)
	}
}
