package member

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"

	"github.com/jarlyyn/herb-go-experimental/user-role"

	"github.com/herb-go/herb/cache"
	"github.com/herb-go/herb/cache-session"
	"github.com/herb-go/herb/middleware"
	"github.com/herb-go/herb/user"

	_ "github.com/herb-go/herb/cache/drivers/freecache"
	"github.com/herb-go/herb/middleware-httprouter"
)

var dataProfileKey = "profile"
var profileData = "herb"

func actionLogin(w http.ResponseWriter, r *http.Request) {
	uid, err := service.Accounts().AccountToUID(*newTestAccount(r.Header.Get("account")))
	if err != nil {
		panic(err)
	}
	err = service.Login(r, uid)
	if err != nil {
		panic(err)
	}
	w.Header().Add("uid", uid)
	_, err = w.Write([]byte("ok"))
	if err != nil {
		panic(err)
	}
}
func actionEcho(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte("ok"))
	if err != nil {
		panic(err)
	}
}

type memberResult struct {
	Accounts user.UserAccounts
	Banned   bool
	Revoke   string
	Role     role.Roles
	Profile  user.Profile
}

func actionMember(w http.ResponseWriter, r *http.Request) {
	var result = memberResult{}
	uid, err := service.IdentifyRequest(r)
	if err != nil {
		panic(err)
	}
	var member = service.GetMembersFromRequest(r)
	accounts, err := member.LoadAccount(uid)
	if err != nil {
		panic(err)
	}
	result.Accounts = accounts[uid]
	banned, err := member.LoadBanned(uid)
	if err != nil {
		panic(err)
	}
	result.Banned = banned[uid]
	revoke, err := member.LoadRevokeTokens(uid)
	if err != nil {
		panic(err)
	}
	result.Revoke = revoke[uid]
	roles, err := member.LoadRoles(uid)
	if err != nil {
		panic(err)
	}
	result.Role = roles[uid]
	profiles, err := member.LoadData(dataProfileKey, uid)
	if err != nil {
		panic(err)
	}
	result.Profile = profiles.(userProfiles)[uid]
	bs, err := json.Marshal(result)
	if err != nil {
		panic(err)
	}
	_, err = w.Write(bs)
	if err != nil {
		panic(err)
	}
}

var service *Service

func initRouter(service *Service, router *httprouter.Router) {
	router.POST("/login").HandleFunc(actionLogin)
	router.POST("/echo").
		Use(
			service.LoginRequiredMiddleware(nil),
			service.BannedMiddleware(),
		).
		HandleFunc(actionEcho)
	router.POST("/role").
		Use(
			service.LoginRequiredMiddleware(nil),
			service.RolesAuthorizeMiddleware("role"),
		).
		HandleFunc(actionEcho)
	router.POST("/member").
		Use(
			service.LoginRequiredMiddleware(nil),
			service.RolesAuthorizeMiddleware(),
		).
		HandleFunc(actionMember)

}
func testService() *Service {

	var store = session.NewClientStore([]byte("12345"), -1)
	config := json.RawMessage("{\"Size\": 10000000}")
	c := cache.New()
	err := c.Open("freecache", config, -1)
	if err != nil {
		panic(err)
	}
	service = NewWithSubCache(store, c)
	service.Install(newTestAccountService())
	service.Install(newTestBannedService())
	service.Install(newTestRevokeService())
	service.Install(newTestPasswordService())
	service.Install(newTestRoleService())
	service.RegisterData(dataProfileKey, *newTestUesrProfiles())
	service.RegisterAccountType("test", user.CaseSensitiveAcountType)
	return service
}
func TestService(t *testing.T) {
	var accountNormalUser = "normalUserAccount"
	var accountNew = "accountNew"
	var password = "password"
	service = testService()
	var app = middleware.New()
	app.Use(service.SessionStore.CookieMiddleware())
	var router = httprouter.New()
	initRouter(service, router)
	app.Handle(router)
	rawUserProfiles = map[string]user.Profile{}
	uid, err := service.Accounts().Register(*newTestAccount(accountNormalUser))
	if err != nil {
		t.Fatal(err)
	}
	rawUserProfiles[uid] = user.Profile{}
	var userprofile = rawUserProfiles[uid]
	userprofile.SetValue(user.ProfileIndexNickname, profileData)
	err = service.Password().UpdatePassword(uid, password)
	if err != nil {
		t.Fatal(err)
	}
	result, err := service.Password().VerifyPassword(uid, password)
	if err != nil {
		t.Fatal(err)
	}
	if result != true {
		t.Error(result)
	}
	var s = httptest.NewServer(app)
	defer s.Close()
	c := s.Client()
	req, err := http.NewRequest("POST", s.URL+"/echo", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 401 {
		t.Error(resp.StatusCode)
	}

	c = s.Client()
	c.Jar, err = cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	req, err = http.NewRequest("POST", s.URL+"/login", nil)
	req.Header.Add("account", accountNormalUser)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Header.Get("uid") != uid {
		t.Error(resp.Header.Get("uid"))
	}

	resp.Body.Close()
	req, err = http.NewRequest("POST", s.URL+"/echo", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	_, err = service.RevokeService.Revoke(uid)
	if err != nil {
		t.Fatal(err)
	}
	req, err = http.NewRequest("POST", s.URL+"/echo", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	revoketoken, err := service.Revoke().Revoke(uid)
	if err != nil {
		t.Fatal(err)
	}
	req, err = http.NewRequest("POST", s.URL+"/echo", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 401 {
		t.Error(resp.StatusCode)
	}
	req, err = http.NewRequest("POST", s.URL+"/login", nil)
	req.Header.Add("account", accountNormalUser)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	req, err = http.NewRequest("POST", s.URL+"/echo", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	err = service.BannedService.Ban(uid, true)
	if err != nil {
		t.Fatal(err)
	}
	req, err = http.NewRequest("POST", s.URL+"/echo", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	err = service.Banned().Ban(uid, true)
	if err != nil {
		t.Fatal(err)
	}
	result, err = service.Password().VerifyPassword(uid, password)
	if err != ErrUserBanned {
		t.Fatal(err)
	}

	req, err = http.NewRequest("POST", s.URL+"/echo", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Error(resp.StatusCode)
	}
	err = service.Banned().Ban(uid, false)
	if err != nil {
		t.Fatal(err)
	}

	req, err = http.NewRequest("POST", s.URL+"/echo", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	req, err = http.NewRequest("POST", s.URL+"/role", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Error(resp.StatusCode)
	}
	var roleservice = service.RoleService.(*testRoleService)
	(*roleservice)[uid] = *role.NewRoles("role", "role2")

	req, err = http.NewRequest("POST", s.URL+"/role", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 403 {
		t.Error(resp.StatusCode)
	}
	err = service.Roles().Clean(uid)
	if err != nil {
		t.Fatal(err)
	}
	req, err = http.NewRequest("POST", s.URL+"/role", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}
	err = service.Accounts().BindAccounts(uid, *newTestAccount(accountNew))
	if err != nil {
		t.Fatal(err)
	}
	req, err = http.NewRequest("POST", s.URL+"/login", nil)
	req.Header.Add("account", accountNew)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	if resp.Header.Get("uid") != uid {
		t.Error(resp.Header.Get("uid"))
	}

	resp.Body.Close()

	req, err = http.NewRequest("POST", s.URL+"/member", nil)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	var memberresult = memberResult{}
	err = json.Unmarshal(content, &memberresult)
	if err != nil {
		t.Fatal(err)
	}
	if len(memberresult.Accounts) != 2 ||
		!memberresult.Accounts.Exists(newTestAccount(accountNormalUser)) ||
		!memberresult.Accounts.Exists(newTestAccount(accountNew)) {
		t.Error(memberresult.Accounts)
	}
	if memberresult.Banned != false {
		t.Error(memberresult.Banned)
	}
	if memberresult.Revoke != revoketoken {
		t.Error(memberresult.Revoke)
	}
	if len(memberresult.Role) != 2 {
		t.Error(memberresult.Role)
	}

	if resp.StatusCode != 200 {
		t.Error(resp.StatusCode)
	}

	err = service.Accounts().UnbindAccounts(uid, *newTestAccount(accountNew))
	if err != nil {
		t.Fatal(err)
	}

	req, err = http.NewRequest("POST", s.URL+"/login", nil)
	req.Header.Add("account", accountNew)
	if err != nil {
		t.Fatal(err)
	}
	resp, err = c.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.Header.Get("uid") != "" {
		t.Error(resp.Header.Get("uid"))
	}
	resp.Body.Close()
}
