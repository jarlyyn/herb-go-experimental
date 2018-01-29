package wechatwork

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/herb-go/herb/fetch"

	auth "github.com/jarlyyn/herb-go-experimental/app-externalauth"
)

const FieldName = "externalauthdriver-wechatwork"
const StateLength = 128
const oauthURL = "https://open.weixin.qq.com/connect/oauth2/authorize"
const qrauthURL = "https://open.work.weixin.qq.com/wwopen/sso/qrConnect"

var DataIndexDepartment = auth.ProfileIndex("WechatWorkDartment")

type Session struct {
	State string
}

func mustHTMLRedirect(w http.ResponseWriter, url string) {
	w.WriteHeader(http.StatusOK)
	html := fmt.Sprintf(`<html><head><meta http-equiv="refresh" content="0; URL='%s'" /></head></html>`, url)
	_, err := w.Write([]byte(html))
	if err != nil {
		panic(err)
	}
}
func authRequestWithAgent(agent *Agent, service *auth.Service, r *http.Request) (*auth.Result, error) {
	var authsession = &Session{}
	q := r.URL.Query()
	var code = q.Get("code")
	if code == "" {
		return nil, nil
	}
	var state = q.Get("state")
	if state == "" {
		return nil, auth.ErrAuthParamsError
	}
	err := service.Auth.Session.Get(r, FieldName, authsession)
	if service.Auth.Session.IsNotFound(err) {
		return nil, nil
	}
	if authsession.State == "" || authsession.State != state {
		return nil, auth.ErrAuthParamsError
	}
	err = service.Auth.Session.Del(r, FieldName)
	if err != nil {
		return nil, err
	}
	info, err := agent.GetUserInfo(code)
	if fetch.CompareApiErrCode(err, ApiErrOauthCodeWrong) {
		return nil, auth.ErrAuthParamsError
	}
	if err != nil {
		return nil, err
	}
	if info == nil {
		return nil, nil
	}
	result := auth.NewResult()
	result.Keyword = service.Keyword
	result.Account = info.UserID
	result.Data.SetValue(auth.ProfileIndexAvatar, info.Avatar)
	result.Data.SetValue(auth.ProfileIndexEmail, info.Email)
	switch info.Gender {
	case ApiResultGenderMale:
		result.Data.SetValue(auth.ProfileIndexGender, auth.ProfileGenderMale)
	case ApiResultGenderFemale:
		result.Data.SetValue(auth.ProfileIndexGender, auth.ProfileGenderFemale)
	}
	result.Data.SetValue(auth.ProfileIndexName, info.Name)
	result.Data.SetValue(auth.ProfileIndexNickname, info.Name)
	for _, v := range info.Department {
		result.Data.AddValue(DataIndexDepartment, strconv.Itoa(v))
	}
	return result, nil
}

type OauthAuthDriver struct {
	agent *Agent
	scope string
}

func NewOauthDriver(agent *Agent, scope string) *OauthAuthDriver {
	return &OauthAuthDriver{
		agent: agent,
		scope: scope,
	}
}

func (d *OauthAuthDriver) ExternalLogin(service *auth.Service, w http.ResponseWriter, r *http.Request) {
	bytes, err := service.Auth.RandToken(StateLength)
	if err != nil {
		panic(err)
	}
	state := string(bytes)
	authsession := Session{
		State: state,
	}
	err = service.Auth.Session.Set(r, FieldName, authsession)
	if err != nil {
		panic(err)
	}
	u, err := url.Parse(oauthURL)
	if err != nil {
		panic(err)
	}
	q := u.Query()
	q.Set("appid", d.agent.CorpID)
	q.Set("agentid", d.agent.AgentID)
	q.Set("scope", d.scope)
	q.Set("state", state)
	q.Set("redirect_uri", service.AuthUrl())
	u.RawQuery = q.Encode()
	u.Fragment = "wechat_redirect"
	mustHTMLRedirect(w, u.String())
}
func (d *OauthAuthDriver) AuthRequest(service *auth.Service, r *http.Request) (*auth.Result, error) {
	return authRequestWithAgent(d.agent, service, r)
}

type QRAuthDriver struct {
	agent *Agent
}

func NewQRAuthDriver(agent *Agent) *QRAuthDriver {
	return &QRAuthDriver{
		agent: agent,
	}
}

func (d *QRAuthDriver) ExternalLogin(service *auth.Service, w http.ResponseWriter, r *http.Request) {
	bytes, err := service.Auth.RandToken(StateLength)
	if err != nil {
		panic(err)
	}
	state := string(bytes)
	authsession := Session{
		State: state,
	}
	err = service.Auth.Session.Set(r, FieldName, authsession)
	if err != nil {
		panic(err)
	}
	u, err := url.Parse(qrauthURL)
	if err != nil {
		panic(err)
	}
	q := u.Query()
	q.Set("appid", d.agent.CorpID)
	q.Set("agentid", d.agent.AgentID)
	q.Set("state", state)
	q.Set("redirect_uri", service.AuthUrl())
	u.RawQuery = q.Encode()
	mustHTMLRedirect(w, u.String())
}
func (d *QRAuthDriver) AuthRequest(service *auth.Service, r *http.Request) (*auth.Result, error) {
	return authRequestWithAgent(d.agent, service, r)
}
